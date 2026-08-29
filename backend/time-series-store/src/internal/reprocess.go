/**
 * @file reprocess.go
 * @brief Obsluha čtení historických dat a synchronizace požadavků pro opětovné vyhodnocení KPI.
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
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

var (
	activeReprocessJobs sync.Map
	deleteWaiters       sync.Map
)

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
	key := makeKey(req.SDTypeUID, req.KPIDefinitionUID)
	activeReprocessJobs.Store(key, req.JobID)
	end := req.To.UTC()
	effectiveFrom, err := c.resolveEffectiveReprocessFrom(req)
	if err != nil {
		return err
	}
	windows := buildTimeWindows(effectiveFrom, end, false)
	needsTerminal := false
	for index, window := range windows {
		flux := c.buildReprocessFluxWindow(req, window.From, window.To)
		hasMoreWindows := index < len(windows)-1
		emittedHasMore, err := c.streamPointsRawFlux(flux, key, req.JobID, req.Batch, onBatch, hasMoreWindows)
		if err != nil {
			return err
		}
		needsTerminal = emittedHasMore
	}
	if needsTerminal {
		return onBatch(nil, false)
	}
	return nil
}

func (c Influx2Client) streamPointsRawFlux(fluxQuery string, jobKey string, jobID string, batchSize int, onBatch func([]sharedModel.TimeSeriesDataPoint, bool) error, forceHasMore bool) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result, err := c.queryApi.Query(ctx, fluxQuery)
	if err != nil {
		return false, err
	}
	pointsBatch := make([]sharedModel.TimeSeriesDataPoint, 0, batchSize)
	emittedHasMore := false
	flush := func(hasMore bool) error {
		if !isActive(jobKey, jobID) {
			cancel()
			return nil
		}
		if len(pointsBatch) == 0 {
			return nil
		}
		emittedHasMore = hasMore
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
			return false, nil
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
				return false, err
			}
		}
	}
	if result.Err() != nil {
		return false, result.Err()
	}
	if err := flush(forceHasMore); err != nil {
		return false, err
	}
	return emittedHasMore, nil
}

func (c Influx2Client) resolveEffectiveReprocessFrom(req sharedModel.TimeSeriesReprocessReadRequest) (time.Time, error) {
	stop := req.To.UTC()
	needGrouping := len(req.SDInstanceUIDs) != 1
	fluxQuery := c.buildReprocessFirstFlux(req, stop, needGrouping)
	result, err := c.queryApi.Query(context.Background(), fluxQuery)
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
		return stop, nil
	}
	return earliest, nil
}

func makeKey(sdType string, kpiUID string) string {
	return fmt.Sprintf("%s|%s", sdType, kpiUID)
}

func isActive(jobKey string, jobID string) bool {
	v, ok := activeReprocessJobs.Load(jobKey)
	if !ok {
		return false
	}
	job, ok := v.(string)
	return ok && job == jobID
}
