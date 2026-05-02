package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func (c Influx2Client) rawPoint(record sharedModel.TimeSeriesRawRecord) (*write.Point, error) {
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
		return nil, err
	}
	fields := map[string]interface{}{
		"payload": payload,
	}
	return influxdb2.NewPoint(measurement, tags, fields, record.EventTime.UTC()), nil
}

func (c Influx2Client) WriteRaw(record sharedModel.TimeSeriesRawRecord) {
	point, err := c.rawPoint(record)
	if err != nil {
		return
	}
	c.writeApi.WritePoint(point)
	c.writeApi.Flush()
}

func (c Influx2Client) WriteRawBatch(records []sharedModel.TimeSeriesRawRecord) {
	if len(records) == 0 {
		return
	}
	points := make([]*write.Point, 0, len(records))
	for _, record := range records {
		point, err := c.rawPoint(record)
		if err != nil {
			continue
		}
		points = append(points, point)
	}
	if len(points) == 0 {
		return
	}
	for _, point := range points {
		c.writeApi.WritePoint(point)
	}
}

func (c Influx2Client) kpiPoint(record sharedModel.TimeSeriesKPIResultRecord) *write.Point {
	if record.JobID != "" {
		if !isActive(makeKey(record.SDTypeUID, record.KPIDefinitionID), record.JobID) {
			return nil
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
	return influxdb2.NewPoint(measurement, tags, fields, record.EventTime.UTC())
}

func (c Influx2Client) WriteKPI(record sharedModel.TimeSeriesKPIResultRecord) {
	point := c.kpiPoint(record)
	if point == nil {
		return
	}
	c.writeApi.WritePoint(point)
	c.writeApi.Flush()
}

func (c Influx2Client) WriteKPIBatch(records []sharedModel.TimeSeriesKPIResultRecord) {
	if len(records) == 0 {
		return
	}
	points := make([]*write.Point, 0, len(records))
	for _, record := range records {
		point := c.kpiPoint(record)
		if point == nil {
			continue
		}
		points = append(points, point)
	}
	if len(points) == 0 {
		return
	}
	for _, point := range points {
		c.writeApi.WritePoint(point)
	}
}

func (c Influx2Client) DeleteKPI(req sharedModel.KPIDeleteResultsRequestISCMessage) error {
	if req.JobID != "" {
		ch := make(chan struct{})
		deleteWaiters.Store(req.JobID, ch)
		defer func() {
			close(ch)
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
