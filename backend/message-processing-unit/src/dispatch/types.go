/**
 * @file types.go
 * @brief Sdílené typy a runtime stav dispatch vrstvy MPU.
 *
 * @author Vojtěch Hubáček
 *
 * @ingroup riot_message_processing_unit
 */
package dispatch

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

type inputProcessingTask struct {
	source string
	tuple  sharedModel.KPIFulfillmentCheckRequestTupleISCMessage
	result chan error
}

type inputSourceStats struct {
	received    atomic.Uint64
	processed   atomic.Uint64
	failed      atomic.Uint64
	dropped     atomic.Uint64
	wouldDrop   atomic.Uint64
	enqueued    atomic.Int64
	paused      atomic.Uint64
	sourcePause atomic.Bool
	trend       inputSourceTrend
}

type inputSourceTrend struct {
	mu                 sync.Mutex
	initialized        bool
	schedulerState     string
	pressureSamples    int
	recoverySamples    int
	dropReadySamples   int
	lastReceived       uint64
	lastProcessed      uint64
	lastInFlight       int64
	shortReceivedRate  float64
	shortProcessedRate float64
	longReceivedRate   float64
	longProcessedRate  float64
	shortBacklogDelta  float64
	longBacklogDelta   float64
}

var (
	ingestConsumerRegistry               sync.Map
	inputSourceStatsMap                  sync.Map
	inputProcessingQueue                 chan inputProcessingTask
	inputWorkersOnce                     sync.Once
	globalInputIntakePaused              atomic.Bool
	inputPauseHighWatermarkPercent       int
	inputPauseLowWatermarkPercent        int
	inputSourcePauseHighWatermarkPercent int
	inputSourcePauseLowWatermarkPercent  int
	inputPauseCheckInterval              time.Duration
	InputTrendSampleInterval             time.Duration
	inputWatchPressureThreshold          float64
	inputDropPressureThreshold           float64
	inputDropBacklogDeltaThreshold       float64
	inputStateEscalationSamples          int
	inputStateRecoverySamples            int
	inputDropperEnabled                  bool
	inputDropperGraceSamples             int
)
