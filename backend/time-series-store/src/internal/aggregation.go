/**
 * @file aggregation.go
 * @brief Sestavení a zpracování agregačních dotazů nad historickými časovými daty.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_time_series_store
 */
package internal

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func extractTags(values map[string]interface{}) map[string]string {
	tags := make(map[string]string)
	for k, v := range values {
		switch k {
		case "_time", "_value", "_field", "_measurement", "_start", "_stop", "_source_timestamp", "result", "table":
			continue
		}
		if str, ok := v.(string); ok {
			tags[k] = str
		}
	}
	return tags
}

func buildAggregateKey(tags map[string]string) string {
	return tags["sdInstanceUID"] + "|" + tags["kpiDefinitionUID"]
}

func cloneTags(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func extractFulfilled(point sharedModel.TimeSeriesDataPoint) (bool, bool) {
	raw, ok := point.Data["fulfilled"]
	if !ok {
		return false, false
	}

	switch v := raw.(type) {
	case bool:
		return v, true
	case string:
		return strings.EqualFold(v, "true"), true
	default:
		return false, false
	}
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func (s *aggregateSeriesState) addDuration(from, to time.Time) {
	if !s.HasLastValue || !to.After(from) {
		return
	}

	d := to.Sub(from)
	if s.LastBool {
		s.TrueWindowTime += d
	} else {
		s.FalseWindowTime += d
	}
	s.LastTime = to
}

func (s *aggregateSeriesState) buildWindowPoint(windowStart time.Time) sharedModel.TimeSeriesDataPoint {
	result := s.LastBool

	if s.TrueWindowTime > s.FalseWindowTime {
		result = true
	} else if s.TrueWindowTime < s.FalseWindowTime {
		result = false
	}

	return sharedModel.TimeSeriesDataPoint{
		Time: windowStart.UTC(),
		Tags: cloneTags(s.Tags),
		Data: map[string]interface{}{
			"fulfilled": result,
		},
	}
}

func (s *aggregateSeriesState) resetWindow() {
	s.TrueWindowTime = 0
	s.FalseWindowTime = 0
	if s.HasLastValue {
		s.HasFullWindow = true
	}
}

func (c Influx2Client) iterateFluxPoints(ctx context.Context, fluxQuery string, plan sharedModel.QueryPlan, isSnapshot bool, onPoint func(sharedModel.TimeSeriesDataPoint) (bool, error)) (bool, error) {
	result, err := c.queryApi.Query(ctx, fluxQuery)
	if err != nil {
		return false, err
	}
	for result.Next() {
		values := result.Record().Values()
		t, ok := values["_time"].(time.Time)
		if !ok {
			continue
		}
		if isSnapshot {
			t = plan.From
		}
		field, _ := values["_field"].(string)
		if field == "" {
			continue
		}
		point := sharedModel.TimeSeriesDataPoint{
			Time: t.UTC(),
			Tags: extractTags(values),
			Data: map[string]interface{}{
				field: values["_value"],
			},
		}
		stop, err := onPoint(point)
		if err != nil {
			return false, err
		}
		if stop {
			return true, nil
		}
	}
	if result.Err() != nil {
		return false, result.Err()
	}
	return false, nil
}

func shouldSkipAggregatedPoint(plan sharedModel.QueryPlan, state *aggregateStreamState, point sharedModel.TimeSeriesDataPoint) bool {
	if !plan.UseCursor || plan.Cursor == nil {
		return false
	}
	if state.cursorPassed {
		return false
	}
	cursorKey := plan.Cursor.SDInstanceUID
	if plan.Cursor.KPIDefinitionUID != nil {
		cursorKey += "|" + *plan.Cursor.KPIDefinitionUID
	}
	pointKey := buildAggregateKey(point.Tags)
	if point.Time.Before(plan.Cursor.Time) {
		return true
	}
	if point.Time.After(plan.Cursor.Time) {
		state.cursorPassed = true
		return false
	}
	if pointKey < cursorKey {
		return true
	}
	if pointKey == cursorKey {
		state.cursorPassed = true
		return true
	}
	state.cursorPassed = true
	return false
}

func (c Influx2Client) flushAggregateWindow(plan sharedModel.QueryPlan, state *aggregateStreamState, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error, hasMoreData bool) (bool, error) {
	for _, key := range state.seriesOrder {
		series := state.series[key]
		if !series.HasLastValue {
			continue
		}
		if !series.HasLastValue {
			continue
		}
		point := series.buildWindowPoint(state.windowStart)
		if shouldSkipAggregatedPoint(plan, state, point) {
			series.resetWindow()
			continue
		}
		state.pointsBatch = append(state.pointsBatch, point)
		state.totalSent++
		series.resetWindow()
		if plan.Limit > 0 && state.totalSent >= plan.Limit {
			if len(state.pointsBatch) > 0 {
				state.needsTerminal = false
				if err := onBatch(state.pointsBatch, false, hasMoreData); err != nil {
					return true, err
				}
				state.pointsBatch = state.pointsBatch[:0]
			}
			return true, nil
		}
		if plan.Batch > 0 && len(state.pointsBatch) >= plan.Batch {
			state.needsTerminal = true
			if err := onBatch(state.pointsBatch, true, hasMoreData); err != nil {
				return false, err
			}
			state.pointsBatch = state.pointsBatch[:0]
		}
	}
	return false, nil
}

func (c Influx2Client) advanceAggregateWindows(plan sharedModel.QueryPlan, state *aggregateStreamState, target time.Time, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error) (bool, error) {
	for state.windowStart.Before(plan.To) {
		windowEnd := minTime(state.windowStart.Add(state.windowSize), plan.To)
		if target.Before(windowEnd) {
			return false, nil
		}
		for _, key := range state.seriesOrder {
			s := state.series[key]
			if !s.HasLastValue {
				continue
			}
			if s.LastTime.Before(windowEnd) {
				s.addDuration(s.LastTime, windowEnd)
			}
		}
		hitLimit, err := c.flushAggregateWindow(plan, state, onBatch, true)
		if err != nil || hitLimit {
			return hitLimit, err
		}
		state.windowStart = windowEnd
	}
	return false, nil
}

func (c Influx2Client) streamReadAggregated(ctx context.Context, plan sharedModel.QueryPlan, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error) error {
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
	effectiveFrom, err := c.resolveEffectivePlanFrom(ctx, plan)
	if err != nil {
		return err
	}
	if !plan.To.After(effectiveFrom) {
		return onBatch(nil, false, false)
	}
	windowSize := time.Duration(*plan.AggregateSeconds) * time.Second
	state := &aggregateStreamState{
		totalSent:   0,
		pointsBatch: make([]sharedModel.TimeSeriesDataPoint, 0, plan.Batch),
		series:      make(map[string]*aggregateSeriesState),
		seriesOrder: make([]string, 0),
		windowStart: effectiveFrom.UTC(),
		windowSize:  windowSize,
	}
	snapshotPlan := plan
	snapshotPlan.From = effectiveFrom
	snapshotFlux := c.buildSnapshotFlux(snapshotPlan)
	_, err = c.iterateFluxPoints(ctx, snapshotFlux, plan, true, func(p sharedModel.TimeSeriesDataPoint) (bool, error) {
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
			LastTime:      effectiveFrom.UTC(),
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
	hitLimit := false
	for _, window := range buildTimeWindows(effectiveFrom, plan.To, false) {
		readFlux := c.buildReadFluxForRange(plan, window.From, window.To)
		chunkHitLimit, err := c.iterateFluxPoints(ctx, readFlux, plan, false, func(p sharedModel.TimeSeriesDataPoint) (bool, error) {
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
		if chunkHitLimit {
			hitLimit = true
			return nil
		}
	}
	hitLimit, err = c.advanceAggregateWindows(plan, state, plan.To.UTC(), onBatch)
	if err != nil {
		return err
	}
	if hitLimit {
		return nil
	}
	if len(state.pointsBatch) > 0 {
		state.needsTerminal = false
		return onBatch(state.pointsBatch, false, false)
	}
	if state.needsTerminal {
		state.needsTerminal = false
		return onBatch(nil, false, false)
	}
	if state.totalSent == 0 {
		return onBatch(nil, false, false)
	}
	return nil
}
