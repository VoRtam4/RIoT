package internal

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

const (
	measurementRaw = "raw_data"
	measurementKPI = "kpi_results"
)

type Influx2Client struct {
	endpoint     string
	organization string
	bucket       string
	client       influxdb2.Client
	writeApi     api.WriteAPIBlocking
	queryApi     api.QueryAPI
}

func NewInflux2Client(endpoint string, token string, organization string, bucket string) Influx2Client {
	log.Printf("[TS] Initializing Influx client | endpoint=%s org=%s bucket=%s", endpoint, organization, bucket)

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

func (c Influx2Client) WriteRaw(record sharedModel.TimeSeriesRawRecord) {
	log.Printf("[TS][RAW] Received raw record | uid=%s type=%s time=%v params=%v",
		record.SDInstanceUID,
		record.SDTypeSpecification,
		record.EventTime,
		record.Parameters,
	)

	if record.EventTime.IsZero() {
		record.EventTime = time.Now().UTC()
	}

	tags := map[string]string{
		"sdInstanceUID": record.SDInstanceUID,
		"sdType":        record.SDTypeSpecification,
	}

	fields := make(map[string]interface{})
	if record.Parameters != nil {
		for k, v := range record.Parameters {
			if v == nil {
				continue
			}
			if k == "uuid" {
				fields["uuid"] = fmt.Sprint(v)
				continue
			}
			fields[k] = v
		}
	}

	if len(fields) == 0 {
		log.Printf("[TS][RAW] Skip write: no fields | uid=%s", record.SDInstanceUID)
		return
	}

	log.Printf("[TS][RAW] Writing point | tags=%v fields=%v", tags, fields)

	point := influxdb2.NewPoint(measurementRaw, tags, fields, record.EventTime.UTC())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.writeApi.WritePoint(ctx, point); err != nil {
		log.Printf("[TS][RAW] WriteRaw failed: %v", err)
		return
	}

	log.Printf("[TS][RAW] Write successful | uid=%s", record.SDInstanceUID)
}

func (c Influx2Client) WriteKPI(record sharedModel.TimeSeriesKPIResultRecord) {
	log.Printf("[TS][KPI] Received KPI result | kpiID=%d uid=%s fulfilled=%v time=%v", record.KPIDefinitionID, record.SDInstanceUID, record.Fulfilled, record.EventTime)
	if record.EventTime.IsZero() {
		record.EventTime = time.Now().UTC()
	}
	fields := map[string]interface{}{
		"fulfilled": record.Fulfilled,
	}
	tags := map[string]string{
		"sdInstanceUID":   record.SDInstanceUID,
		"kpiDefinitionID": fmt.Sprintf("%d", record.KPIDefinitionID),
	}
	if record.SDTypeSpecification != "" {
		tags["sdType"] = record.SDTypeSpecification
	}
	log.Printf("[TS][KPI] Writing KPI point | tags=%v fields=%v", tags, fields)
	point := influxdb2.NewPoint(measurementKPI, tags, fields, record.EventTime.UTC())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.writeApi.WritePoint(ctx, point); err != nil {
		log.Printf("[TS][KPI] WriteKPI failed: %v", err)
		return
	}
	log.Printf("[TS][KPI] Write successful | kpiID=%d uid=%s", record.KPIDefinitionID, record.SDInstanceUID)
}

func (c Influx2Client) Query(req sharedModel.TimeSeriesReadRequest) sharedUtils.Result[[]sharedModel.TimeSeriesDataPoint] {
	log.Printf("[TS][QUERY] Incoming request | type=%v uid=%v sdType=%v kpiID=%v from=%v to=%v normalize=%v limit=%v batch=%v", req.Type, req.SDInstanceUID, req.SDTypeSpecification, req.KPIDefinitionID, req.From, req.To, req.Normalize, req.Limit, req.Batch)
	points, err := c.collectPoints(req)
	if err != nil {
		return sharedUtils.NewFailureResult[[]sharedModel.TimeSeriesDataPoint](err)
	}
	points = applyLimitPreservingSameTimestamp(points, req)
	log.Printf("[TS][QUERY] Final collected points = %d", len(points))
	return sharedUtils.NewSuccessResult(points)
}

func (c Influx2Client) StreamQueryBatches(req sharedModel.TimeSeriesReadRequest, onBatch func(points []sharedModel.TimeSeriesDataPoint, hasMore bool) error) error {
	// Normalize a Limit nejdou korektně řešit plně streamingově bez dočtení / doskládání.
	// Pro tyto případy fallback na collect mód a jedna odpověď.
	if shouldNormalize(req) || hasLimit(req) {
		result := c.Query(req)
		if result.IsFailure() {
			return result.GetError()
		}
		return onBatch(result.GetPayload(), false)
	}
	batchSize := resolveBatchSize(req)
	return c.streamPointsNoBoundary(req, batchSize, onBatch)
}

func (c Influx2Client) collectPoints(req sharedModel.TimeSeriesReadRequest) ([]sharedModel.TimeSeriesDataPoint, error) {
	normalize := shouldNormalize(req)
	mainReq := req
	mainReq.Normalize = boolPtr(false)
	mainPoints := make([]sharedModel.TimeSeriesDataPoint, 0)
	err := c.streamPointsNoBoundary(mainReq, 0, func(points []sharedModel.TimeSeriesDataPoint, hasMore bool) error {
		mainPoints = append(mainPoints, points...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !normalize || req.From == nil {
		return mainPoints, nil
	}
	boundaryPoints, err := c.queryBoundaryPoints(req)
	if err != nil {
		return nil, err
	}
	if len(boundaryPoints) == 0 {
		return mainPoints, nil
	}
	from := req.From.UTC()
	for i := range boundaryPoints {
		boundaryPoints[i].Time = from
	}
	sortBoundaryPoints(boundaryPoints, req.Type, false)
	desc := false
	if req.SortDesc != nil {
		desc = *req.SortDesc
	}
	if desc {
		return append(mainPoints, boundaryPoints...), nil
	}
	return append(boundaryPoints, mainPoints...), nil
}

func (c Influx2Client) queryBoundaryPoints(req sharedModel.TimeSeriesReadRequest) ([]sharedModel.TimeSeriesDataPoint, error) {
	if req.From == nil {
		return nil, nil
	}
	stop := req.From.UTC().Add(-time.Nanosecond)
	if stop.Before(time.Unix(0, 0).UTC()) {
		return nil, nil
	}
	boundaryReq := req
	boundaryReq.From = nil
	boundaryReq.To = &stop
	boundaryReq.Normalize = boolPtr(false)
	boundaryReq.Limit = nil
	boundaryReq.Batch = nil
	desc := true
	boundaryReq.SortDesc = &desc
	allPoints := make([]sharedModel.TimeSeriesDataPoint, 0)
	err := c.streamPointsNoBoundary(boundaryReq, 0, func(points []sharedModel.TimeSeriesDataPoint, hasMore bool) error {
		allPoints = append(allPoints, points...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(allPoints) == 0 {
		return nil, nil
	}
	lastTime := allPoints[0].Time.UTC()
	result := make([]sharedModel.TimeSeriesDataPoint, 0)
	for _, p := range allPoints {
		if !p.Time.UTC().Equal(lastTime) {
			break
		}
		result = append(result, p)
	}
	return result, nil
}

func (c Influx2Client) streamPointsNoBoundary(req sharedModel.TimeSeriesReadRequest, batchSize int, onBatch func(points []sharedModel.TimeSeriesDataPoint, hasMore bool) error) error {
	fluxQuery, err := c.buildFastFluxQuery(req)
	if err != nil {
		return err
	}
	log.Printf("[TS][QUERY] Flux query:\n%s", fluxQuery)
	result, err := c.queryApi.Query(context.Background(), fluxQuery)
	if err != nil {
		return err
	}
	pointsBatch := make([]sharedModel.TimeSeriesDataPoint, 0)
	var current *sharedModel.TimeSeriesDataPoint
	currentKey := ""
	emittedSomething := false
	flushBatch := func(hasMore bool) error {
		if len(pointsBatch) == 0 {
			return nil
		}
		emittedSomething = true
		err := onBatch(pointsBatch, hasMore)
		if err != nil {
			return err
		}
		pointsBatch = make([]sharedModel.TimeSeriesDataPoint, 0)
		return nil
	}
	flushCurrentPoint := func() error {
		if current == nil {
			return nil
		}
		pointsBatch = append(pointsBatch, *current)
		current = nil
		currentKey = ""

		if batchSize > 0 && len(pointsBatch) >= batchSize {
			return flushBatch(true)
		}
		return nil
	}
	for result.Next() {
		values := result.Record().Values()
		t, ok := values["_time"].(time.Time)
		if !ok {
			continue
		}
		pointKey := buildPointKeyFromRow(values, req.Type)
		if current == nil || currentKey != pointKey {
			if err := flushCurrentPoint(); err != nil {
				return err
			}
			p := sharedModel.TimeSeriesDataPoint{
				Time: t.UTC(),
				Tags: map[string]string{},
				Data: map[string]interface{}{},
			}
			copyTagIfPresent(values, p.Tags, "sdInstanceUID")
			copyTagIfPresent(values, p.Tags, "sdType")
			copyTagIfPresent(values, p.Tags, "kpiDefinitionID")
			current = &p
			currentKey = pointKey
		}
		fieldName, ok := values["_field"].(string)
		if !ok || fieldName == "" {
			continue
		}
		fieldValue, ok := values["_value"]
		if !ok || fieldValue == nil {
			continue
		}
		current.Data[fieldName] = fieldValue
	}
	if result.Err() != nil {
		return result.Err()
	}
	if err := flushCurrentPoint(); err != nil {
		return err
	}
	if len(pointsBatch) > 0 {
		return flushBatch(false)
	}
	if !emittedSomething {
		return onBatch([]sharedModel.TimeSeriesDataPoint{}, false)
	}
	return nil
}

func (c Influx2Client) buildFastFluxQuery(req sharedModel.TimeSeriesReadRequest) (string, error) {
	measurement := measurementRaw
	if req.Type == sharedModel.TimeSeriesTypeKPIResult {
		measurement = measurementKPI
	}
	rangePart, err := buildRange(req.From, req.To)
	if err != nil {
		return "", err
	}
	filterPart := buildFilters(req, measurement)
	aggregationPart := buildAggregationPart(req)
	sortPart := buildFastSortPart(req.Type, req.SortDesc)
	query := fmt.Sprintf(`
from(bucket: "%s")
%s
%s
`, c.bucket, rangePart, filterPart)
	if aggregationPart != "" {
		query += aggregationPart + "\n"
	}
	if sortPart != "" {
		query += sortPart + "\n"
	}
	return query, nil
}

func shouldNormalize(req sharedModel.TimeSeriesReadRequest) bool {
	if req.Normalize != nil {
		return *req.Normalize
	}
	return true
}

func hasLimit(req sharedModel.TimeSeriesReadRequest) bool {
	return req.Limit != nil && *req.Limit > 0
}

func resolveBatchSize(req sharedModel.TimeSeriesReadRequest) int {
	if req.Batch == nil || *req.Batch == 0 {
		return 0
	}
	return int(*req.Batch)
}

func buildRange(from *time.Time, to *time.Time) (string, error) {
	if from != nil && to != nil && from.UTC().After(to.UTC()) {
		return "", fmt.Errorf("invalid time interval: from (%s) is after to (%s)", from.UTC().Format(time.RFC3339Nano), to.UTC().Format(time.RFC3339Nano))
	}
	if from == nil && to == nil {
		return `|> range(start:0)`, nil
	}
	if from != nil && to == nil {
		return fmt.Sprintf(`|> range(start:%s)`, from.UTC().Format(time.RFC3339Nano)), nil
	}
	if from == nil && to != nil {
		stopInclusive := to.UTC().Add(time.Nanosecond)
		return fmt.Sprintf(`|> range(start:0, stop:%s)`, stopInclusive.Format(time.RFC3339Nano)), nil
	}
	stopInclusive := to.UTC().Add(time.Nanosecond)
	return fmt.Sprintf(`|> range(start:%s, stop:%s)`, from.UTC().Format(time.RFC3339Nano), stopInclusive.Format(time.RFC3339Nano)), nil
}

func buildFilters(req sharedModel.TimeSeriesReadRequest, measurement string) string {
	filter := fmt.Sprintf(`|> filter(fn: (r) => r["_measurement"] == "%s")`, measurement)
	if req.SDInstanceUID != nil && *req.SDInstanceUID != "" {
		filter += fmt.Sprintf(` |> filter(fn: (r) => r["sdInstanceUID"] == "%s")`, escapeFluxString(*req.SDInstanceUID))
	}
	if req.SDTypeSpecification != nil && *req.SDTypeSpecification != "" {
		filter += fmt.Sprintf(` |> filter(fn: (r) => r["sdType"] == "%s")`, escapeFluxString(*req.SDTypeSpecification))
	}
	if req.KPIDefinitionID != nil {
		filter += fmt.Sprintf(` |> filter(fn: (r) => r["kpiDefinitionID"] == "%d")`, *req.KPIDefinitionID)
	}
	return filter
}

func buildAggregationPart(req sharedModel.TimeSeriesReadRequest) string {
	if req.AggregateMinutes == nil || *req.AggregateMinutes <= 0 {
		return ""
	}
	return fmt.Sprintf(`|> aggregateWindow(every:%dm, fn:last, createEmpty:false)`, *req.AggregateMinutes)
}

func buildFastSortPart(tsType sharedModel.TimeSeriesType, sortDesc *bool) string {
	desc := false
	if sortDesc != nil {
		desc = *sortDesc
	}
	columns := []string{`"_time"`, `"sdInstanceUID"`}
	if tsType == sharedModel.TimeSeriesTypeKPIResult {
		columns = []string{`"_time"`, `"sdInstanceUID"`, `"kpiDefinitionID"`}
	}
	columns = append(columns, `"_field"`)
	return fmt.Sprintf(`|> sort(columns:[%s], desc:%v)`, strings.Join(columns, ","), desc)
}

func applyLimitPreservingSameTimestamp(points []sharedModel.TimeSeriesDataPoint, req sharedModel.TimeSeriesReadRequest) []sharedModel.TimeSeriesDataPoint {
	if req.Limit == nil || *req.Limit <= 0 {
		return points
	}
	limit := *req.Limit
	if len(points) <= limit {
		return points
	}
	cutoffTime := points[limit-1].Time
	i := limit
	for i < len(points) && points[i].Time.Equal(cutoffTime) {
		i++
	}
	return points[:i]
}

func buildPointKeyFromRow(values map[string]interface{}, tsType sharedModel.TimeSeriesType) string {
	t, _ := values["_time"].(time.Time)
	sdInstanceUID := safeStringFromValues(values, "sdInstanceUID")
	sdType := safeStringFromValues(values, "sdType")
	kpiDefinitionID := safeStringFromValues(values, "kpiDefinitionID")
	if tsType == sharedModel.TimeSeriesTypeKPIResult {
		return fmt.Sprintf("%d|%s|%s|%s", t.UTC().UnixNano(), sdInstanceUID, sdType, kpiDefinitionID)
	}
	return fmt.Sprintf("%d|%s|%s", t.UTC().UnixNano(), sdInstanceUID, sdType)
}

func safeStringFromValues(values map[string]interface{}, key string) string {
	v, ok := values[key]
	if !ok || v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func sortBoundaryPoints(points []sharedModel.TimeSeriesDataPoint, tsType sharedModel.TimeSeriesType, desc bool) {
	sort.Slice(points, func(i, j int) bool {
		return pointLess(points[i], points[j], tsType, desc)
	})
}

func pointLess(a sharedModel.TimeSeriesDataPoint, b sharedModel.TimeSeriesDataPoint, tsType sharedModel.TimeSeriesType, desc bool) bool {
	if !a.Time.Equal(b.Time) {
		if desc {
			return a.Time.After(b.Time)
		}
		return a.Time.Before(b.Time)
	}
	aUID := a.Tags["sdInstanceUID"]
	bUID := b.Tags["sdInstanceUID"]
	if aUID != bUID {
		if desc {
			return aUID > bUID
		}
		return aUID < bUID
	}
	if tsType == sharedModel.TimeSeriesTypeKPIResult {
		aKPI := a.Tags["kpiDefinitionID"]
		bKPI := b.Tags["kpiDefinitionID"]
		if aKPI != bKPI {
			if desc {
				return aKPI > bKPI
			}
			return aKPI < bKPI
		}
	}
	return false
}

func copyTagIfPresent(values map[string]interface{}, tags map[string]string, key string) {
	if v, ok := values[key]; ok && v != nil {
		tags[key] = fmt.Sprint(v)
	}
}

func escapeFluxString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

func boolPtr(v bool) *bool {
	return &v
}

func (c Influx2Client) Close() {
	log.Printf("[TS] Closing Influx client")
	c.client.Close()
}

func (c Influx2Client) DeleteKPI(kpiDefinitionID uint32) {
	log.Printf("[TS][KPI] Delete request | kpiID=%d", kpiDefinitionID)
	predicate := fmt.Sprintf(`_measurement="%s" AND kpiDefinitionID="%d"`, measurementKPI, kpiDefinitionID)
	start := time.Unix(0, 0)
	stop := time.Now().UTC()
	deleteAPI := c.client.DeleteAPI()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := deleteAPI.DeleteWithName(ctx, c.organization, c.bucket, start, stop, predicate)
	if err != nil {
		log.Printf("[TS][KPI] DeleteKPI failed: %v", err)
		return
	}
	log.Printf("[TS][KPI] Delete successful | kpiID=%d", kpiDefinitionID)
}
