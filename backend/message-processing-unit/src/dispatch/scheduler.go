/**
 * @file scheduler.go
 * @brief Pause/resume, trendová klasifikace a dropper rozhodování dispatch vrstvy MPU.
 *
 * @author Vojtěch Hubáček
 *
 * @ingroup riot_message_processing_unit
 */
package dispatch

import (
	"log"
	"time"
)

func updateInputSourceSchedulerState(source string, trend *inputSourceTrend) string {
	shortPressure := trend.shortReceivedRate - trend.shortProcessedRate
	longPressure := trend.longReceivedRate - trend.longProcessedRate
	overWatch := shortPressure >= inputWatchPressureThreshold || trend.shortBacklogDelta > 0
	overDropCandidate := longPressure >= inputDropPressureThreshold && trend.longBacklogDelta >= inputDropBacklogDeltaThreshold
	recovering := shortPressure <= 0 && longPressure <= 0 && trend.shortBacklogDelta <= 0

	if overDropCandidate {
		trend.pressureSamples++
		trend.recoverySamples = 0
	} else if overWatch {
		if trend.pressureSamples < inputStateEscalationSamples {
			trend.pressureSamples++
		}
		trend.recoverySamples = 0
	} else if recovering {
		trend.recoverySamples++
		if trend.pressureSamples > 0 {
			trend.pressureSamples--
		}
	} else {
		trend.recoverySamples = 0
	}

	previousState := trend.schedulerState
	switch {
	case overDropCandidate && trend.pressureSamples >= inputStateEscalationSamples:
		trend.schedulerState = "dropCandidate"
		trend.dropReadySamples++
	case overWatch && trend.pressureSamples >= inputStateEscalationSamples:
		trend.schedulerState = "watch"
		trend.dropReadySamples = 0
	case trend.schedulerState == "dropCandidate" && !overDropCandidate:
		trend.schedulerState = "cooldown"
		trend.dropReadySamples = 0
	case trend.schedulerState == "watch" && !overWatch:
		trend.schedulerState = "cooldown"
		trend.dropReadySamples = 0
	case trend.schedulerState == "cooldown" && trend.recoverySamples >= inputStateRecoverySamples:
		trend.schedulerState = "normal"
		trend.dropReadySamples = 0
	case trend.schedulerState == "":
		trend.schedulerState = "normal"
		trend.dropReadySamples = 0
	}

	if previousState != trend.schedulerState {
		log.Printf(
			"[MPU][INPUT][SCHEDULER] Source state changed | source=%s state=%s previous=%s shortPressure=%.2f longPressure=%.2f shortBacklogDelta=%.2f longBacklogDelta=%.2f pressureSamples=%d recoverySamples=%d dropReadySamples=%d",
			source,
			trend.schedulerState,
			previousState,
			shortPressure,
			longPressure,
			trend.shortBacklogDelta,
			trend.longBacklogDelta,
			trend.pressureSamples,
			trend.recoverySamples,
			trend.dropReadySamples,
		)
	}
	return trend.schedulerState
}

func isInputSourceDropReady(stats *inputSourceStats) bool {
	stats.trend.mu.Lock()
	defer stats.trend.mu.Unlock()
	return stats.trend.schedulerState == "dropCandidate" && stats.trend.dropReadySamples >= inputDropperGraceSamples
}

func updateGlobalInputPauseState() (bool, int, int, int) {
	if inputProcessingQueue == nil || cap(inputProcessingQueue) == 0 {
		globalInputIntakePaused.Store(false)
		return false, 0, 0, 0
	}
	queueLen := len(inputProcessingQueue)
	queueCap := cap(inputProcessingQueue)
	loadPercent := queueLen * 100 / queueCap
	wasPaused := globalInputIntakePaused.Load()
	if wasPaused {
		if loadPercent <= inputPauseLowWatermarkPercent {
			globalInputIntakePaused.Store(false)
			log.Printf("[MPU][INPUT][SCHEDULER] Resuming global intake | internalQueue=%d/%d load=%d%%", queueLen, queueCap, loadPercent)
		}
	} else if loadPercent >= inputPauseHighWatermarkPercent {
		globalInputIntakePaused.Store(true)
		log.Printf("[MPU][INPUT][SCHEDULER] Pausing global intake | internalQueue=%d/%d load=%d%%", queueLen, queueCap, loadPercent)
	}
	return globalInputIntakePaused.Load(), loadPercent, queueLen, queueCap
}

func updateInputSourcePauseState(source string, stats *inputSourceStats) bool {
	if inputProcessingQueue == nil || cap(inputProcessingQueue) == 0 {
		stats.sourcePause.Store(false)
		return false
	}
	sourceInFlight := stats.enqueued.Load()
	queueCap := cap(inputProcessingQueue)
	sourceLoadPercent := int(sourceInFlight) * 100 / queueCap
	wasPaused := stats.sourcePause.Load()
	if wasPaused {
		if sourceLoadPercent <= inputSourcePauseLowWatermarkPercent {
			stats.sourcePause.Store(false)
			log.Printf("[MPU][INPUT][SCHEDULER] Resuming source intake | source=%s inFlight=%d capacity=%d load=%d%%", source, sourceInFlight, queueCap, sourceLoadPercent)
		}
	} else if sourceLoadPercent >= inputSourcePauseHighWatermarkPercent {
		stats.sourcePause.Store(true)
		log.Printf("[MPU][INPUT][SCHEDULER] Pausing source intake | source=%s inFlight=%d capacity=%d load=%d%%", source, sourceInFlight, queueCap, sourceLoadPercent)
	}
	return stats.sourcePause.Load()
}

func refreshInputSourcePauseStates() {
	inputSourceStatsMap.Range(func(key any, value any) bool {
		updateInputSourcePauseState(key.(string), value.(*inputSourceStats))
		return true
	})
}

func waitForInputCapacity(source string) {
	stats := getInputSourceStats(source)
	for {
		globalPaused, _, _, _ := updateGlobalInputPauseState()
		sourcePaused := updateInputSourcePauseState(source, stats)
		if !globalPaused && !sourcePaused {
			return
		}
		stats.paused.Add(1)
		time.Sleep(inputPauseCheckInterval)
	}
}
