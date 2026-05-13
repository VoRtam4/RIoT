/**
 * @file read.go
 * @brief Streamované čtení historických dat z InfluxDB a jejich převod do interního modelu.
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
	"log"
	"sort"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func (c Influx2Client) StreamRead(ctx context.Context, plan sharedModel.QueryPlan, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error) error {
	if plan.UseAggregation {
		return c.streamReadAggregated(ctx, plan, onBatch)
	}
	state := &streamState{
		totalSent:    0,
		pointsBatch:  make([]sharedModel.TimeSeriesDataPoint, 0, plan.Batch),
		cursorPassed: false,
		hadAnyData:   false,
	}
	finalFlush := func() error {
		if len(state.pointsBatch) > 0 {
			state.needsTerminal = false
			return onBatch(state.pointsBatch, false, false)
		}
		if state.needsTerminal {
			state.needsTerminal = false
			return onBatch(nil, false, false)
		}
		if state.totalSent == 0 {
			if state.hadAnyData {
				log.Printf("[TS][WARNING] data existed but nothing sent")
			}
			return onBatch(nil, false, false)
		}
		return nil
	}
	if !plan.UseSort {
		hitLimit, err := c.streamReadRangeChunks(ctx, plan, state, onBatch)
		if err != nil {
			return err
		}
		if hitLimit && state.totalSent != 0 {
			return nil
		}
		return finalFlush()
	}
	if !plan.SortDesc {
		if plan.NeedInitial {
			snapshotFlux := c.buildSnapshotFlux(plan)
			hitLimit, err := c.streamRowsFlux(ctx, snapshotFlux, plan, state, onBatch, true)
			if err != nil {
				return err
			}
			if hitLimit && state.totalSent != 0 {
				return nil
			}
		}
		hitLimit, err := c.streamReadRangeChunks(ctx, plan, state, onBatch)
		if err != nil {
			return err
		}
		if hitLimit && state.totalSent != 0 {
			return nil
		}
		return finalFlush()
	}
	hitLimit, err := c.streamReadRangeChunks(ctx, plan, state, onBatch)
	if err != nil {
		return err
	}
	if hitLimit && state.totalSent != 0 {
		return nil
	}
	if !plan.From.IsZero() {
		snapshotFlux := c.buildSnapshotFlux(plan)
		hitLimit, err := c.streamRowsFlux(ctx, snapshotFlux, plan, state, onBatch, true)
		if err != nil {
			return err
		}
		if hitLimit && state.totalSent != 0 {
			return nil
		}
	}
	return finalFlush()
}

func (c Influx2Client) streamReadRangeChunks(ctx context.Context, plan sharedModel.QueryPlan, state *streamState, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error) (bool, error) {
	if !plan.To.After(plan.From) {
		return false, nil
	}

	effectiveFrom, err := c.resolveEffectivePlanFrom(ctx, plan)
	if err != nil {
		return false, err
	}

	windows := buildTimeWindows(effectiveFrom, plan.To, plan.UseSort && plan.SortDesc)
	for _, window := range windows {
		readFlux := c.buildReadFluxForRange(plan, window.From, window.To)
		hitLimit, err := c.streamRowsFlux(ctx, readFlux, plan, state, onBatch, false)
		if err != nil || hitLimit {
			return hitLimit, err
		}
	}
	return false, nil
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
	distinctValues := make(map[string]struct{})
	effectiveFrom, err := c.resolveEffectivePlanFrom(context.Background(), plan)
	if err != nil {
		return nil, err
	}
	for _, window := range buildTimeWindows(effectiveFrom, plan.To, false) {
		fluxQuery := c.buildDistinctTagValuesFluxForRange(plan, req.Tag, window.From, window.To)
		result, err := c.queryApi.Query(context.Background(), fluxQuery)
		if err != nil {
			return nil, err
		}

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
	}

	values := make([]string, 0, len(distinctValues))
	for value := range distinctValues {
		values = append(values, value)
	}
	sort.Strings(values)
	return values, nil
}

func (c Influx2Client) resolveEffectivePlanFrom(ctx context.Context, plan sharedModel.QueryPlan) (time.Time, error) {
	if !plan.From.IsZero() || !plan.To.After(plan.From) {
		return plan.From, nil
	}

	fluxQuery := c.buildFirstPointFlux(plan, plan.To)
	result, err := c.queryApi.Query(ctx, fluxQuery)
	if err != nil {
		return time.Time{}, err
	}

	var earliest time.Time
	found := false
	for result.Next() {
		values := result.Record().Values()
		t, ok := values["_time"].(time.Time)
		if !ok {
			continue
		}
		if !found || t.Before(earliest) {
			earliest = t.UTC()
			found = true
		}
	}
	if result.Err() != nil {
		return time.Time{}, result.Err()
	}
	if !found {
		return plan.To, nil
	}
	return earliest, nil
}

func (c Influx2Client) streamRowsFlux(ctx context.Context, fluxQuery string, plan sharedModel.QueryPlan, state *streamState, onBatch func([]sharedModel.TimeSeriesDataPoint, bool, bool) error, isSnapshot bool) (bool, error) {
	result, err := c.queryApi.Query(ctx, fluxQuery)
	if err != nil {
		return false, err
	}
	flush := func(hasMoreBatches bool, hasMoreData bool) error {
		if len(state.pointsBatch) == 0 {
			return nil
		}
		state.needsTerminal = hasMoreBatches
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
