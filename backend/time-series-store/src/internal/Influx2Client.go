package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	sharedModel "github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
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
	if record.EventTime.IsZero() {
		record.EventTime = time.Now().UTC()
	}

	tags := map[string]string{
		"sdInstanceUID": record.SDInstanceUID,
		"sdType":        record.SDTypeSpecification,
	}

	if record.Source != "" {
		tags["source"] = record.Source
	}

	fields := make(map[string]interface{})

	if record.Parameters != nil {
		if raw, ok := record.Parameters["uuid"]; ok && raw != nil {
			tags["uuid"] = fmt.Sprint(raw)
		}

		for k, v := range record.Parameters {
			if v == nil {
				continue
			}
			if k == "uuid" {
				continue
			}
			fields[k] = v
		}
	}

	if len(fields) == 0 {
		return
	}

	point := influxdb2.NewPoint(
		measurementRaw,
		tags,
		fields,
		record.EventTime.UTC(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.writeApi.WritePoint(ctx, point); err != nil {
		log.Printf("WriteRaw failed: %v", err)
	}
}

func (c Influx2Client) WriteKPI(record sharedModel.TimeSeriesKPIResultRecord) {
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

	if record.Source != "" {
		tags["source"] = record.Source
	}

	point := influxdb2.NewPoint(
		measurementKPI,
		tags,
		fields,
		record.EventTime.UTC(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.writeApi.WritePoint(ctx, point); err != nil {
		log.Printf("WriteKPI failed: %v", err)
	}
}

func (c Influx2Client) Query(req sharedModel.TimeSeriesReadRequest) sharedUtils.Result[[]sharedModel.TimeSeriesDataPoint] {

	measurement := measurementRaw
	if req.Type == sharedModel.TimeSeriesTypeKPIResult {
		measurement = measurementKPI
	}

	rangePart := buildRange(req)
	filterPart := buildFilters(req, measurement)
	aggregationPart := buildAggregation(req)

	query := fmt.Sprintf(`
from(bucket: "%s")
%s
%s
%s
|> pivot(rowKey:["_time"], columnKey:["_field"], valueColumn:"_value")
`, c.bucket, rangePart, filterPart, aggregationPart)

	result, err := c.queryApi.Query(context.Background(), query)
	if err != nil {
		return sharedUtils.NewFailureResult[[]sharedModel.TimeSeriesDataPoint](err)
	}

	if result.Err() != nil {
		return sharedUtils.NewFailureResult[[]sharedModel.TimeSeriesDataPoint](result.Err())
	}

	points := make([]sharedModel.TimeSeriesDataPoint, 0)

	for result.Next() {

		values := result.Record().Values()

		t, ok := values["_time"].(time.Time)
		if !ok {
			continue
		}

		point := sharedModel.TimeSeriesDataPoint{
			Time: t,
			Tags: map[string]string{},
			Data: map[string]interface{}{},
		}

		copyTagIfPresent(values, point.Tags, "sdInstanceUID")
		copyTagIfPresent(values, point.Tags, "sdType")
		copyTagIfPresent(values, point.Tags, "source")
		copyTagIfPresent(values, point.Tags, "kpiDefinitionID")

		for k, v := range values {

			if k == "_time" || k == "_measurement" || k == "result" || k == "table" {
				continue
			}

			if k == "sdInstanceUID" || k == "sdType" || k == "source" || k == "kpiDefinitionID" {
				continue
			}

			if v == nil {
				continue
			}

			point.Data[k] = v
		}

		points = append(points, point)
	}

	return sharedUtils.NewSuccessResult(points)
}

func (c Influx2Client) Close() {
	c.client.Close()
}

func buildRange(req sharedModel.TimeSeriesReadRequest) string {

	if req.From == nil && req.To == nil {
		return `|> range(start: -30d)`
	}

	if req.From != nil && req.To == nil {
		return fmt.Sprintf(`|> range(start: %s)`, req.From.UTC().Format(time.RFC3339))
	}

	if req.From == nil && req.To != nil {
		start := req.To.UTC().AddDate(0, 0, -30)
		return fmt.Sprintf(`|> range(start: %s, stop: %s)`,
			start.Format(time.RFC3339),
			req.To.UTC().Format(time.RFC3339))
	}

	return fmt.Sprintf(`|> range(start: %s, stop: %s)`,
		req.From.UTC().Format(time.RFC3339),
		req.To.UTC().Format(time.RFC3339))
}

func buildFilters(req sharedModel.TimeSeriesReadRequest, measurement string) string {

	filter := fmt.Sprintf(`|> filter(fn: (r) => r["_measurement"] == "%s")`, measurement)

	if req.SDInstanceUID != nil && *req.SDInstanceUID != "" {
		filter += fmt.Sprintf(` |> filter(fn: (r) => r["sdInstanceUID"] == "%s")`,
			escapeFluxString(*req.SDInstanceUID))
	}

	if req.SDTypeSpecification != nil && *req.SDTypeSpecification != "" {
		filter += fmt.Sprintf(` |> filter(fn: (r) => r["sdType"] == "%s")`,
			escapeFluxString(*req.SDTypeSpecification))
	}

	if req.Source != nil && *req.Source != "" {
		filter += fmt.Sprintf(` |> filter(fn: (r) => r["source"] == "%s")`,
			escapeFluxString(*req.Source))
	}

	if req.KPIDefinitionID != nil {
		filter += fmt.Sprintf(` |> filter(fn: (r) => r["kpiDefinitionID"] == "%d")`,
			*req.KPIDefinitionID)
	}

	return filter
}

func buildAggregation(req sharedModel.TimeSeriesReadRequest) string {

	if req.AggregateMinutes == nil || *req.AggregateMinutes <= 0 {
		return ""
	}

	return fmt.Sprintf(`|> aggregateWindow(every: %dm, fn: mean, createEmpty: false)`,
		*req.AggregateMinutes)
}

func copyTagIfPresent(values map[string]interface{}, tags map[string]string, key string) {

	if v, ok := values[key]; ok && v != nil {
		if s, ok := v.(string); ok && s != "" {
			tags[key] = s
		}
	}
}

func escapeFluxString(s string) string {

	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)

	return s
}
