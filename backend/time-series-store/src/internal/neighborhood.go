package internal

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

const (
	neighborhoodLocalPointLimit = 32
	neighborhoodDefaultGap      = time.Second
)

type timelineSummary struct {
	First time.Time
	Last  time.Time
	Count int
}

func (c Influx2Client) GetRecordNeighborhood(req sharedModel.TimeSeriesRecordNeighborhoodRequest) (sharedModel.TimeSeriesRecordNeighborhoodResponse, error) {
	if req.SDTypeUID == "" || req.SDInstanceUID == "" || req.EventTime.IsZero() {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, fmt.Errorf("missing sdTypeUID, sdInstanceUID or eventTime")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	epsilon := neighborhoodDefaultGap
	if summary, err := c.rawTimelineSummary(ctx, req); err == nil && summary.Count > 1 {
		roughGap := summary.Last.Sub(summary.First) / time.Duration(summary.Count-1)
		if roughGap > 0 {
			epsilon = roughGap
		}
	}
	if medianGap, err := c.localMedianGap(ctx, req, epsilon); err == nil && medianGap > 0 {
		epsilon = medianGap
	}
	if epsilon <= 0 {
		epsilon = neighborhoodDefaultGap
	}

	response, err := c.rawNeighborhoodWithin(ctx, req, epsilon*2)
	if err != nil {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, err
	}
	if response.Previous != nil {
		response.PreviousKPIState, _ = c.kpiStateAt(ctx, req, response.Previous.Time)
	}
	return response, nil
}

func (c Influx2Client) rawTimelineSummary(ctx context.Context, req sharedModel.TimeSeriesRecordNeighborhoodRequest) (timelineSummary, error) {
	flux := c.buildRawTimelineSummaryFlux(req)
	result, err := c.queryApi.Query(ctx, flux)
	if err != nil {
		return timelineSummary{}, err
	}
	summary := timelineSummary{}
	for result.Next() {
		values := result.Record().Values()
		metric, _ := values["metric"].(string)
		switch metric {
		case "first", "last":
			t, ok := values["_time"].(time.Time)
			if !ok {
				continue
			}
			if metric == "first" {
				summary.First = t.UTC()
			} else {
				summary.Last = t.UTC()
			}
		case "count":
			count, ok := intValue(values["_value"])
			if ok {
				summary.Count = count
			}
		}
	}
	if result.Err() != nil {
		return timelineSummary{}, result.Err()
	}
	return summary, nil
}

func (c Influx2Client) localMedianGap(ctx context.Context, req sharedModel.TimeSeriesRecordNeighborhoodRequest, initialRange time.Duration) (time.Duration, error) {
	if initialRange <= 0 {
		initialRange = neighborhoodDefaultGap
	}
	window := initialRange
	for attempt := 0; attempt < 3; attempt++ {
		points, err := c.queryRawPoints(ctx, req, req.EventTime.Add(-window), req.EventTime.Add(window), false, neighborhoodLocalPointLimit)
		if err != nil {
			return 0, err
		}
		if len(points) >= 2 {
			gaps := make([]time.Duration, 0, len(points)-1)
			for i := 1; i < len(points); i++ {
				gap := points[i].Time.Sub(points[i-1].Time)
				if gap > 0 {
					gaps = append(gaps, gap)
				}
			}
			if len(gaps) > 0 {
				sort.Slice(gaps, func(i, j int) bool { return gaps[i] < gaps[j] })
				return gaps[len(gaps)/2], nil
			}
		}
		window *= 2
	}
	return 0, nil
}

func intValue(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case uint64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		parsed, err := strconv.Atoi(v)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func (c Influx2Client) rawNeighborhoodWithin(ctx context.Context, req sharedModel.TimeSeriesRecordNeighborhoodRequest, window time.Duration) (sharedModel.TimeSeriesRecordNeighborhoodResponse, error) {
	points, err := c.queryRawPoints(ctx, req, req.EventTime.Add(-window), req.EventTime.Add(window).Add(time.Nanosecond), false, neighborhoodLocalPointLimit*2)
	if err != nil {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, err
	}
	response := sharedModel.TimeSeriesRecordNeighborhoodResponse{}
	for i := range points {
		point := points[i]
		switch {
		case point.Time.Equal(req.EventTime):
			response.SameTimestamp = &point
		case point.Time.Before(req.EventTime):
			if response.Previous == nil || point.Time.After(response.Previous.Time) {
				response.Previous = &point
			}
		case point.Time.After(req.EventTime):
			if response.Next == nil || point.Time.Before(response.Next.Time) {
				response.Next = &point
			}
		}
	}
	return response, nil
}

func (c Influx2Client) queryRawPoints(ctx context.Context, req sharedModel.TimeSeriesRecordNeighborhoodRequest, from time.Time, to time.Time, desc bool, limit int) ([]sharedModel.TimeSeriesDataPoint, error) {
	if !to.After(from) {
		return nil, nil
	}
	flux := c.buildRawPointsFlux(req, from, to, desc, limit)
	result, err := c.queryApi.Query(ctx, flux)
	if err != nil {
		return nil, err
	}
	points := make([]sharedModel.TimeSeriesDataPoint, 0)
	for result.Next() {
		values := result.Record().Values()
		t, ok := values["_time"].(time.Time)
		if !ok {
			continue
		}
		payload, ok := values["_value"].(string)
		if !ok {
			continue
		}
		data, err := sharedUtils.DecodeCBOR(payload)
		if err != nil {
			continue
		}
		points = append(points, sharedModel.TimeSeriesDataPoint{
			Time: t.UTC(),
			Tags: extractTags(values),
			Data: data,
		})
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	sort.Slice(points, func(i, j int) bool {
		if desc {
			return points[i].Time.After(points[j].Time)
		}
		return points[i].Time.Before(points[j].Time)
	})
	return points, nil
}

func (c Influx2Client) kpiStateAt(ctx context.Context, req sharedModel.TimeSeriesRecordNeighborhoodRequest, timestamp time.Time) (map[string]bool, error) {
	flux := c.buildKPIStateAtFlux(req, timestamp)
	result, err := c.queryApi.Query(ctx, flux)
	if err != nil {
		return nil, err
	}
	state := make(map[string]bool)
	for result.Next() {
		values := result.Record().Values()
		kpiUID, _ := values["kpiDefinitionUID"].(string)
		if kpiUID == "" {
			continue
		}
		switch value := values["_value"].(type) {
		case bool:
			state[kpiUID] = value
		}
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	return state, nil
}
