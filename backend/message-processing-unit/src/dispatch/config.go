/**
 * @file config.go
 * @brief Načítání konfigurace dispatch vrstvy MPU z proměnných prostředí.
 *
 * @author Vojtěch Hubáček
 *
 * @ingroup riot_message_processing_unit
 */
package dispatch

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func getPositiveInt(envName string, defaultValue int, logPrefix string) int {
	raw := sharedUtils.GetEnvironmentVariableValue(envName).GetPayloadOrDefault(strconv.Itoa(defaultValue))
	value, convertErr := strconv.Atoi(raw)
	if convertErr != nil || value < 1 {
		log.Printf("%s Invalid %s=%q, using %d", logPrefix, envName, raw, defaultValue)
		return defaultValue
	}
	return value
}

func getPercentInt(envName string, defaultValue int, logPrefix string) int {
	raw := sharedUtils.GetEnvironmentVariableValue(envName).GetPayloadOrDefault(strconv.Itoa(defaultValue))
	value, convertErr := strconv.Atoi(raw)
	if convertErr != nil || value < 1 || value > 100 {
		log.Printf("%s Invalid %s=%q, using %d", logPrefix, envName, raw, defaultValue)
		return defaultValue
	}
	return value
}

func getFloat(envName string, defaultValue float64, logPrefix string) float64 {
	raw := sharedUtils.GetEnvironmentVariableValue(envName).GetPayloadOrDefault(fmt.Sprintf("%.2f", defaultValue))
	value, convertErr := strconv.ParseFloat(raw, 64)
	if convertErr != nil || value < 0 {
		log.Printf("%s Invalid %s=%q, using %.2f", logPrefix, envName, raw, defaultValue)
		return defaultValue
	}
	return value
}

func getBool(envName string, defaultValue bool, logPrefix string) bool {
	rawDefault := "false"
	if defaultValue {
		rawDefault = "true"
	}
	raw := sharedUtils.GetEnvironmentVariableValue(envName).GetPayloadOrDefault(rawDefault)
	value, convertErr := strconv.ParseBool(raw)
	if convertErr != nil {
		log.Printf("%s Invalid %s=%q, using %t", logPrefix, envName, raw, defaultValue)
		return defaultValue
	}
	return value
}

func configureInputScheduler(logPrefix string) {
	inputPauseHighWatermarkPercent = getPercentInt("MPU_INPUT_PAUSE_HIGH_WATERMARK_PERCENT", 80, logPrefix)
	inputPauseLowWatermarkPercent = getPercentInt("MPU_INPUT_RESUME_LOW_WATERMARK_PERCENT", 50, logPrefix)
	if inputPauseLowWatermarkPercent >= inputPauseHighWatermarkPercent {
		log.Printf(
			"%s Invalid input pause hysteresis low=%d high=%d, using low=50 high=80",
			logPrefix,
			inputPauseLowWatermarkPercent,
			inputPauseHighWatermarkPercent,
		)
		inputPauseLowWatermarkPercent = 50
		inputPauseHighWatermarkPercent = 80
	}
	inputSourcePauseHighWatermarkPercent = getPercentInt("MPU_INPUT_SOURCE_PAUSE_HIGH_WATERMARK_PERCENT", 50, logPrefix)
	inputSourcePauseLowWatermarkPercent = getPercentInt("MPU_INPUT_SOURCE_RESUME_LOW_WATERMARK_PERCENT", 25, logPrefix)
	if inputSourcePauseLowWatermarkPercent >= inputSourcePauseHighWatermarkPercent {
		log.Printf(
			"%s Invalid per-source input pause hysteresis low=%d high=%d, using low=25 high=50",
			logPrefix,
			inputSourcePauseLowWatermarkPercent,
			inputSourcePauseHighWatermarkPercent,
		)
		inputSourcePauseLowWatermarkPercent = 25
		inputSourcePauseHighWatermarkPercent = 50
	}
	pauseCheckIntervalMs := getPositiveInt("MPU_INPUT_PAUSE_CHECK_INTERVAL_MS", 200, logPrefix)
	inputPauseCheckInterval = time.Duration(pauseCheckIntervalMs) * time.Millisecond
	trendSampleIntervalMs := getPositiveInt("MPU_INPUT_TREND_SAMPLE_INTERVAL_MS", 5000, logPrefix)
	InputTrendSampleInterval = time.Duration(trendSampleIntervalMs) * time.Millisecond
	inputWatchPressureThreshold = getFloat("MPU_INPUT_WATCH_PRESSURE_THRESHOLD", 0.1, logPrefix)
	inputDropPressureThreshold = getFloat("MPU_INPUT_DROP_CANDIDATE_PRESSURE_THRESHOLD", 0.5, logPrefix)
	inputDropBacklogDeltaThreshold = getFloat("MPU_INPUT_DROP_CANDIDATE_BACKLOG_DELTA_THRESHOLD", 1, logPrefix)
	inputStateEscalationSamples = getPositiveInt("MPU_INPUT_STATE_ESCALATION_SAMPLES", 3, logPrefix)
	inputStateRecoverySamples = getPositiveInt("MPU_INPUT_STATE_RECOVERY_SAMPLES", 3, logPrefix)
	inputDropperEnabled = getBool("MPU_INPUT_DROPPER_ENABLED", false, logPrefix)
	inputDropperGraceSamples = getPositiveInt("MPU_INPUT_DROPPER_GRACE_SAMPLES", 3, logPrefix)
	log.Printf(
		"%s Scheduler configured | globalPauseHigh=%d%% globalResumeLow=%d%% sourcePauseHigh=%d%% sourceResumeLow=%d%% pauseCheck=%s trendSample=%s watchPressure=%.2f dropPressure=%.2f dropBacklogDelta=%.2f escalationSamples=%d recoverySamples=%d dropperEnabled=%t dropperGraceSamples=%d",
		logPrefix,
		inputPauseHighWatermarkPercent,
		inputPauseLowWatermarkPercent,
		inputSourcePauseHighWatermarkPercent,
		inputSourcePauseLowWatermarkPercent,
		inputPauseCheckInterval,
		InputTrendSampleInterval,
		inputWatchPressureThreshold,
		inputDropPressureThreshold,
		inputDropBacklogDeltaThreshold,
		inputStateEscalationSamples,
		inputStateRecoverySamples,
		inputDropperEnabled,
		inputDropperGraceSamples,
	)
}
