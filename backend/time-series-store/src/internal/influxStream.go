package internal

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

var (
	activeReprocessJobs sync.Map
	deleteWaiters       sync.Map
)

func (c Influx2Client) streamPointsRawFlux(fluxQuery string, jobKey string, jobID string, batchSize int, onBatch func([]sharedModel.TimeSeriesDataPoint, bool) error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result, err := c.queryApi.Query(ctx, fluxQuery)
	if err != nil {
		return err
	}
	pointsBatch := make([]sharedModel.TimeSeriesDataPoint, 0, batchSize)
	flush := func(hasMore bool) error {
		if !isActive(jobKey, jobID) {
			cancel()
			return nil
		}
		if len(pointsBatch) == 0 {
			return nil
		}
		if err := onBatch(pointsBatch, hasMore); err != nil {
			return err
		}
		pointsBatch = pointsBatch[:0]
		return nil
	}
	hasNext := result.Next()
	for hasNext {
		if !isActive(jobKey, jobID) {
			cancel()
			return nil
		}
		rec := result.Record()
		values := rec.Values()
		t, ok := values["_time"].(time.Time)
		if !ok {
			hasNext = result.Next()
			continue
		}
		payload, ok := values["_value"].(string)
		if !ok {
			hasNext = result.Next()
			continue
		}
		data, err := sharedUtils.DecodeCBOR(payload)
		if err != nil {
			hasNext = result.Next()
			continue
		}
		point := sharedModel.TimeSeriesDataPoint{
			Time: t.UTC(),
			Tags: make(map[string]string),
			Data: data,
		}
		for k, v := range extractTags(values) {
			point.Tags[k] = fmt.Sprint(v)
		}
		pointsBatch = append(pointsBatch, point)
		hasNext = result.Next()
		if batchSize > 0 && len(pointsBatch) >= batchSize {
			if err := flush(hasNext); err != nil {
				return err
			}
		}
	}
	if result.Err() != nil {
		return result.Err()
	}
	return flush(false)
}

func (c Influx2Client) streamRowsFlux(fluxQuery string, plan sharedModel.QueryPlan, state *streamState, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error, isSnapshot bool) (bool, error) {
	result, err := c.queryApi.Query(context.Background(), fluxQuery)
	if err != nil {
		return false, err
	}
	flush := func(hasMoreBatches bool, hasMoreData bool) error {
		if len(state.pointsBatch) == 0 {
			return nil
		}
		err := onBatch(state.pointsBatch, hasMoreBatches, hasMoreData)
		state.pointsBatch = state.pointsBatch[:0]
		return err
	}
	buildCursorKey := func(tags map[string]string) string {
		key := tags["sdInstanceUID"]
		if plan.IsKPI {
			key += "|" + tags["kpiDefinitionID"]
		}
		return key
	}
	cursorKey := ""
	if plan.UseCursor && plan.Cursor != nil {
		cursorKey = plan.Cursor.SDInstanceUID
		if plan.IsKPI && plan.Cursor.KPIDefinitionID != nil {
			cursorKey += "|" + fmt.Sprintf("%d", *plan.Cursor.KPIDefinitionID)
		}
	}
	shouldSkip := func(p sharedModel.TimeSeriesDataPoint) bool {
		if !plan.UseCursor || plan.Cursor == nil {
			return false
		}
		if state.cursorPassed {
			return false
		}
		if plan.SortDesc {
			if p.Time.After(plan.Cursor.Time) {
				return true
			}
		} else {
			if p.Time.Before(plan.Cursor.Time) {
				return true
			}
		}
		if !p.Time.Equal(plan.Cursor.Time) {
			state.cursorPassed = true
			return false
		}
		if p.Time.Equal(plan.Cursor.Time) && buildCursorKey(p.Tags) == cursorKey {
			state.cursorPassed = true
			return true
		}
		return true
	}
	hasNext := result.Next()
	for hasNext {
		values := result.Record().Values()
		t, ok := values["_time"].(time.Time)
		if !ok {
			hasNext = result.Next()
			continue
		}
		if isSnapshot {
			t = plan.From
		}
		var point sharedModel.TimeSeriesDataPoint
		if plan.IsKPI {
			field, _ := values["_field"].(string)
			if field == "" {
				hasNext = result.Next()
				continue
			}
			point = sharedModel.TimeSeriesDataPoint{
				Time: t.UTC(),
				Tags: extractTags(values),
				Data: map[string]interface{}{
					field: values["_value"],
				},
			}
		} else {
			payload, ok := values["_value"].(string)
			if !ok {
				hasNext = result.Next()
				continue
			}
			data, err := sharedUtils.DecodeCBOR(payload)
			if err != nil {
				hasNext = result.Next()
				continue
			}
			point = sharedModel.TimeSeriesDataPoint{
				Time: t.UTC(),
				Tags: extractTags(values),
				Data: data,
			}
		}
		state.hadAnyData = true
		if shouldSkip(point) {
			hasNext = result.Next()
			continue
		}
		state.pointsBatch = append(state.pointsBatch, point)
		state.totalSent++
		hasNext = result.Next()
		if plan.Limit > 0 && state.totalSent >= plan.Limit {
			return true, flush(false, true)
		}
		if plan.Batch > 0 && len(state.pointsBatch) >= plan.Batch && hasNext {
			if err := flush(true, true); err != nil {
				return false, err
			}
		}
	}
	if result.Err() != nil {
		return false, result.Err()
	}
	return false, nil
}

func (c Influx2Client) streamReadAggregated(plan sharedModel.QueryPlan, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error) error {
	if !plan.IsKPI {
		return fmt.Errorf("aggregation only supported for KPI")
	}
	if plan.SortDesc {
		return fmt.Errorf("descending aggregation not implemented yet")
	}
	if plan.AggregateSeconds == nil || *plan.AggregateSeconds <= 0 {
		return fmt.Errorf("invalid aggregate seconds")
	}
	if !plan.To.After(plan.From) {
		return onBatch(nil, false, false)
	}
	windowSize := time.Duration(*plan.AggregateSeconds) * time.Second
	state := &aggregateStreamState{
		totalSent:   0,
		pointsBatch: make([]sharedModel.TimeSeriesDataPoint, 0, plan.Batch),
		series:      make(map[string]*aggregateSeriesState),
		seriesOrder: make([]string, 0),
		windowStart: plan.From.UTC(),
		windowSize:  windowSize,
	}
	snapshotFlux := c.buildSnapshotFlux(plan)
	_, err := c.iterateFluxPoints(snapshotFlux, plan, true, func(p sharedModel.TimeSeriesDataPoint) (bool, error) {
		key := buildAggregateKey(p.Tags)
		val, ok := extractFulfilled(p)
		if !ok {
			return false, nil
		}
		if _, exists := state.series[key]; exists {
			return false, nil
		}
		state.series[key] = &aggregateSeriesState{
			Key:           key,
			Tags:          cloneTags(p.Tags),
			LastBool:      val,
			LastTime:      plan.From.UTC(),
			HasLastValue:  true,
			HasFullWindow: true,
		}
		state.seriesOrder = append(state.seriesOrder, key)
		return false, nil
	})
	if err != nil {
		return err
	}
	sort.Strings(state.seriesOrder)
	readFlux := c.buildReadFlux(plan)
	hitLimit, err := c.iterateFluxPoints(readFlux, plan, false, func(p sharedModel.TimeSeriesDataPoint) (bool, error) {
		t := p.Time.UTC()
		hitLimit, err := c.advanceAggregateWindows(plan, state, t, onBatch)
		if err != nil || hitLimit {
			return hitLimit, err
		}
		key := buildAggregateKey(p.Tags)
		val, ok := extractFulfilled(p)
		if !ok {
			return false, nil
		}
		series, exists := state.series[key]
		if !exists {
			series = &aggregateSeriesState{
				Key:           key,
				Tags:          cloneTags(p.Tags),
				LastBool:      val,
				LastTime:      t,
				HasLastValue:  true,
				HasFullWindow: false,
			}
			state.series[key] = series
			state.seriesOrder = append(state.seriesOrder, key)
			sort.Strings(state.seriesOrder)
			return false, nil
		}
		if series.HasLastValue && t.After(series.LastTime) {
			series.addDuration(series.LastTime, t)
		}
		series.LastBool = val
		series.LastTime = t
		return false, nil
	})
	if err != nil {
		return err
	}
	if hitLimit {
		return nil
	}
	hitLimit, err = c.advanceAggregateWindows(plan, state, plan.To.UTC(), onBatch)
	if err != nil {
		return err
	}
	if hitLimit {
		return nil
	}
	if len(state.pointsBatch) > 0 {
		return onBatch(state.pointsBatch, false, false)
	}
	if state.totalSent == 0 {
		return onBatch(nil, false, false)
	}
	return nil
}
