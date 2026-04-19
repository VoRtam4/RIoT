package internal

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInflux2Client(endpoint string, token string, organization string, bucket string) Influx2Client {
	client := influxdb2.NewClientWithOptions(endpoint, token, influxdb2.DefaultOptions().SetBatchSize(20))
	return Influx2Client{
		endpoint:     endpoint,
		organization: organization,
		bucket:       bucket,
		client:       client,
		writeApi:     client.WriteAPIBlocking(organization, bucket),
		queryApi:     client.QueryAPI(organization),
	}
}

func (c Influx2Client) Close() {
	c.client.Close()
}

func (c Influx2Client) WriteRaw(record sharedModel.TimeSeriesRawRecord) {
	log.Printf("%v", record)
	if record.EventTime.IsZero() {
		record.EventTime = time.Now().UTC()
	}
	measurement := fmt.Sprintf("%s_%s", string(sharedModel.TimeSeriesTypeRaw), record.SDTypeUID)
	tags := map[string]string{
		"sdInstanceUID": record.SDInstanceUID,
	}
	for k, v := range record.Tags {
		tags[k] = v
	}
	payload, err := sharedUtils.EncodeCBOR(record.Fields)
	if err != nil {
		return
	}
	fields := map[string]interface{}{
		"payload": payload,
	}
	point := influxdb2.NewPoint(measurement, tags, fields, record.EventTime.UTC())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = c.writeApi.WritePoint(ctx, point)
}

func (c Influx2Client) WriteKPI(record sharedModel.TimeSeriesKPIResultRecord) {
	if record.JobID != "" {
		if !isActive(makeKey(record.SDTypeUID, record.KPIDefinitionID), record.JobID) {
			return
		}
	}
	if record.EventTime.IsZero() {
		record.EventTime = time.Now().UTC()
	}
	measurement := fmt.Sprintf("%s_%s", string(sharedModel.TimeSeriesTypeKPIResult), record.SDTypeUID)
	fields := map[string]interface{}{"fulfilled": record.Fulfilled}
	tags := map[string]string{"sdInstanceUID": record.SDInstanceUID, "kpiDefinitionID": fmt.Sprintf("%d", record.KPIDefinitionID)}
	for k, v := range record.Tags {
		tags[k] = v
	}
	point := influxdb2.NewPoint(measurement, tags, fields, record.EventTime.UTC())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.writeApi.WritePoint(ctx, point); err != nil {
		return
	}
}

func (c Influx2Client) DeleteKPI(req sharedModel.KPIDeleteResultsRequestISCMessage) error {
	if req.JobID != "" {
		ch := make(chan struct{})
		deleteWaiters.Store(req.JobID, ch)
		defer func() {
			close(ch)
			deleteWaiters.Delete(req.JobID)
		}()
	}
	activeReprocessJobs.Delete(makeKey(req.SDTypeUID, req.KPIDefinitionID))
	predicate := fmt.Sprintf(`_measurement="%s_%s" AND kpiDefinitionID="%d"`, string(sharedModel.TimeSeriesTypeKPIResult), req.SDTypeUID, req.KPIDefinitionID)
	start := time.Unix(0, 0)
	stop := time.Now().UTC()
	deleteAPI := c.client.DeleteAPI()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := deleteAPI.DeleteWithName(ctx, c.organization, c.bucket, start, stop, predicate)
	if err != nil {
		return nil
	}
	return nil
}

func (c Influx2Client) StreamReprocess(req sharedModel.TimeSeriesReprocessReadRequest, onBatch func([]sharedModel.TimeSeriesDataPoint, bool) error) error {
	if req.JobID != "" && req.Wait {
		var ch chan struct{}
		for {
			if val, ok := deleteWaiters.Load(req.JobID); ok {
				ch = val.(chan struct{})
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		<-ch
	}
	flux := c.buildReprocessFlux(req)
	key := makeKey(req.SDTypeUID, req.KPIDefinitionID)
	activeReprocessJobs.Store(key, req.JobID)
	err := c.streamPointsRawFlux(flux, key, req.JobID, req.Batch, onBatch)
	return err
}

func (c Influx2Client) StreamRead(plan sharedModel.QueryPlan, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error) error {
	if plan.UseAggregation {
		return c.streamReadAggregated(plan, onBatch)
	}
	state := &streamState{
		totalSent:    0,
		pointsBatch:  make([]sharedModel.TimeSeriesDataPoint, 0, plan.Batch),
		cursorPassed: false,
		hadAnyData:   false,
	}
	finalFlush := func() error {
		if len(state.pointsBatch) > 0 {
			return onBatch(state.pointsBatch, false, false)
		}
		if state.totalSent == 0 {
			if state.hadAnyData {
				log.Printf("[TS][WARNING] data existed but nothing sent")
			}
			return onBatch(nil, false, false)
		}
		return nil
	}
	if !plan.SortDesc {
		if plan.NeedInitial {
			snapshotFlux := c.buildSnapshotFlux(plan)
			hitLimit, err := c.streamRowsFlux(snapshotFlux, plan, state, onBatch, true)
			if err != nil {
				return err
			}
			if hitLimit && state.totalSent != 0 {
				return nil
			}
		}
		if plan.To.After(plan.From) {
			readFlux := c.buildReadFlux(plan)
			hitLimit, err := c.streamRowsFlux(readFlux, plan, state, onBatch, false)
			if err != nil {
				return err
			}
			if hitLimit && state.totalSent != 0 {
				return nil
			}
		}
		return finalFlush()
	}
	if plan.To.After(plan.From) {
		readFlux := c.buildReadFlux(plan)
		hitLimit, err := c.streamRowsFlux(readFlux, plan, state, onBatch, false)
		if err != nil {
			return err
		}
		if hitLimit && state.totalSent != 0 {
			return nil
		}
	}
	if !plan.From.IsZero() {
		snapshotFlux := c.buildSnapshotFlux(plan)
		hitLimit, err := c.streamRowsFlux(snapshotFlux, plan, state, onBatch, true)
		if err != nil {
			return err
		}
		if hitLimit && state.totalSent != 0 {
			return nil
		}
	}
	return finalFlush()
}

func (c Influx2Client) DistinctTagValues(req sharedModel.TimeSeriesDistinctTagValuesRequest) ([]string, error) {
	readReq := sharedModel.TimeSeriesReadRequest{
		Type:             req.Type,
		SDTypeUID:        req.SDTypeUID,
		SDInstanceUIDs:   req.SDInstanceUIDs,
		KPIDefinitionIDs: req.KPIDefinitionIDs,
		From:             req.From,
		To:               req.To,
		Filters:          req.Filters,
	}
	plan := BuildQueryPlan(readReq)
	fluxQuery := c.buildDistinctTagValuesFlux(plan, req.Tag)

	result, err := c.queryApi.Query(context.Background(), fluxQuery)
	if err != nil {
		return nil, err
	}

	distinctValues := make(map[string]struct{})
	for result.Next() {
		value := result.Record().ValueByKey(req.Tag)
		if value == nil {
			value = result.Record().Value()
		}
		if value == nil {
			continue
		}
		switch typed := value.(type) {
		case string:
			if typed != "" {
				distinctValues[typed] = struct{}{}
			}
		case []byte:
			if len(typed) > 0 {
				distinctValues[string(typed)] = struct{}{}
			}
		default:
			stringified := fmt.Sprint(typed)
			if stringified != "" && stringified != "<nil>" {
				distinctValues[stringified] = struct{}{}
			}
		}
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	values := make([]string, 0, len(distinctValues))
	for value := range distinctValues {
		values = append(values, value)
	}
	sort.Strings(values)
	return values, nil
}
