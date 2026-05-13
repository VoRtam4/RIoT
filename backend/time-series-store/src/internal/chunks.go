/**
 * @file chunks.go
 * @brief Pomocná logika pro dělení časových intervalů dotazů na stabilní zpracovatelné úseky.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_time_series_store
 */
package internal

import "time"

const queryChunkDuration = 25 * time.Hour

type timeWindow struct {
	From time.Time
	To   time.Time
}

func normalizeChunkBoundary(t time.Time) time.Time {
	if t.IsZero() {
		return time.Unix(0, 0).UTC()
	}
	return t.UTC()
}

func buildTimeWindows(from, to time.Time, descending bool) []timeWindow {
	if from.IsZero() {
		if !to.After(from) {
			return nil
		}
		return []timeWindow{{From: from, To: to.UTC()}}
	}

	start := normalizeChunkBoundary(from)
	end := normalizeChunkBoundary(to)
	if !end.After(start) {
		return nil
	}

	windows := make([]timeWindow, 0, int(end.Sub(start)/queryChunkDuration)+1)
	if descending {
		cursor := end
		for cursor.After(start) {
			windowStart := cursor.Add(-queryChunkDuration)
			if windowStart.Before(start) {
				windowStart = start
			}
			windows = append(windows, timeWindow{
				From: windowStart,
				To:   cursor,
			})
			cursor = windowStart
		}
		return windows
	}

	cursor := start
	for cursor.Before(end) {
		windowEnd := cursor.Add(queryChunkDuration)
		if windowEnd.After(end) {
			windowEnd = end
		}
		windows = append(windows, timeWindow{
			From: cursor,
			To:   windowEnd,
		})
		cursor = windowEnd
	}
	return windows
}
