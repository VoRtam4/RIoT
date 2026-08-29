/**
 * @file stats.go
 * @brief Statistiky, trendy a periodické logování dispatch vrstvy MPU.
 *
 * @author Vojtěch Hubáček
 *
 * @ingroup riot_message_processing_unit
 */
package dispatch

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

func getInputSourceStats(source string) *inputSourceStats {
	rawStats, _ := inputSourceStatsMap.LoadOrStore(source, &inputSourceStats{})
	return rawStats.(*inputSourceStats)
}

func inputStatsSnapshot() []string {
	lines := make([]string, 0)
	inputSourceStatsMap.Range(func(key any, value any) bool {
		source := key.(string)
		stats := value.(*inputSourceStats)
		lines = append(lines, fmt.Sprintf(
			"source=%s received=%d processed=%d failed=%d dropped=%d wouldDrop=%d inFlight=%d sourcePaused=%t paused=%d",
			source,
			stats.received.Load(),
			stats.processed.Load(),
			stats.failed.Load(),
			stats.dropped.Load(),
			stats.wouldDrop.Load(),
			stats.enqueued.Load(),
			stats.sourcePause.Load(),
			stats.paused.Load(),
		))
		return true
	})
	sort.Strings(lines)
	return lines
}

func ewma(previous float64, current float64, alpha float64) float64 {
	return alpha*current + (1-alpha)*previous
}

func inputTrendSnapshot() []string {
	lines := make([]string, 0)
	intervalSeconds := InputTrendSampleInterval.Seconds()
	if intervalSeconds <= 0 {
		intervalSeconds = 1
	}
	inputSourceStatsMap.Range(func(key any, value any) bool {
		source := key.(string)
		stats := value.(*inputSourceStats)
		received := stats.received.Load()
		processed := stats.processed.Load()
		inFlight := stats.enqueued.Load()

		stats.trend.mu.Lock()
		defer stats.trend.mu.Unlock()

		currentReceivedRate := float64(received-stats.trend.lastReceived) / intervalSeconds
		currentProcessedRate := float64(processed-stats.trend.lastProcessed) / intervalSeconds
		currentBacklogDelta := float64(inFlight - stats.trend.lastInFlight)

		if !stats.trend.initialized {
			stats.trend.schedulerState = "normal"
			stats.trend.shortReceivedRate = currentReceivedRate
			stats.trend.shortProcessedRate = currentProcessedRate
			stats.trend.longReceivedRate = currentReceivedRate
			stats.trend.longProcessedRate = currentProcessedRate
			stats.trend.shortBacklogDelta = currentBacklogDelta
			stats.trend.longBacklogDelta = currentBacklogDelta
			stats.trend.initialized = true
		} else {
			stats.trend.shortReceivedRate = ewma(stats.trend.shortReceivedRate, currentReceivedRate, 0.5)
			stats.trend.shortProcessedRate = ewma(stats.trend.shortProcessedRate, currentProcessedRate, 0.5)
			stats.trend.longReceivedRate = ewma(stats.trend.longReceivedRate, currentReceivedRate, 0.15)
			stats.trend.longProcessedRate = ewma(stats.trend.longProcessedRate, currentProcessedRate, 0.15)
			stats.trend.shortBacklogDelta = ewma(stats.trend.shortBacklogDelta, currentBacklogDelta, 0.5)
			stats.trend.longBacklogDelta = ewma(stats.trend.longBacklogDelta, currentBacklogDelta, 0.15)
		}

		stats.trend.lastReceived = received
		stats.trend.lastProcessed = processed
		stats.trend.lastInFlight = inFlight
		state := updateInputSourceSchedulerState(source, &stats.trend)

		lines = append(lines, fmt.Sprintf(
			"source=%s state=%s recv/s=%.2f/%.2f proc/s=%.2f/%.2f pressure=%.2f/%.2f backlogDelta=%.2f/%.2f dropReadySamples=%d inFlight=%d",
			source,
			state,
			stats.trend.shortReceivedRate,
			stats.trend.longReceivedRate,
			stats.trend.shortProcessedRate,
			stats.trend.longProcessedRate,
			stats.trend.shortReceivedRate-stats.trend.shortProcessedRate,
			stats.trend.longReceivedRate-stats.trend.longProcessedRate,
			stats.trend.shortBacklogDelta,
			stats.trend.longBacklogDelta,
			stats.trend.dropReadySamples,
			inFlight,
		))
		return true
	})
	sort.Strings(lines)
	return lines
}

func LogInputStatsPeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		queueLen := 0
		queueCap := 0
		if inputProcessingQueue != nil {
			queueLen = len(inputProcessingQueue)
			queueCap = cap(inputProcessingQueue)
		}
		globalPaused, _, _, _ := updateGlobalInputPauseState()
		refreshInputSourcePauseStates()
		statsLines := inputStatsSnapshot()
		if len(statsLines) == 0 {
			log.Printf("[MPU][INPUT][STATS] internalQueue=%d/%d globalPaused=%t sources=0", queueLen, queueCap, globalPaused)
			continue
		}
		log.Printf("[MPU][INPUT][STATS] internalQueue=%d/%d globalPaused=%t sources=%d | %s", queueLen, queueCap, globalPaused, len(statsLines), strings.Join(statsLines, " | "))
	}
}

func LogInputTrendsPeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		trendLines := inputTrendSnapshot()
		if len(trendLines) == 0 {
			log.Printf("[MPU][INPUT][TREND] interval=%s sources=0", interval)
			continue
		}
		log.Printf("[MPU][INPUT][TREND] interval=%s sources=%d | %s", interval, len(trendLines), strings.Join(trendLines, " | "))
	}
}
