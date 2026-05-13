/**
 * @file tupleProcessing.go
 * @brief Dávkové zpracování vstupních tuple zpráv a orchestrace navazujícího raw a KPI toku.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_message_processing_unit
 */
package processing

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

type normalizedKPICheckRequest struct {
	Message   sharedModel.KPIFulfillmentCheckRequestISCMessage
	Params    map[string]interface{}
	EventTime time.Time
}

var (
	inputConcurrency     atomic.Int32
	instanceProcessLocks sync.Map
)

func SetInputConcurrency(workerCount int) {
	if workerCount < 1 {
		workerCount = 1
	}
	inputConcurrency.Store(int32(workerCount))
}

func lockInstanceForInputProcessing(uid string) func() {
	if inputConcurrency.Load() <= 1 || uid == "" {
		return func() {}
	}
	rawLock, _ := instanceProcessLocks.LoadOrStore(uid, &sync.Mutex{})
	lock := rawLock.(*sync.Mutex)
	lock.Lock()
	return lock.Unlock
}

func ProcessKPIFulfillmentCheckRequestTuple(rabbitMQClient rabbitmq.Client, tuple sharedModel.KPIFulfillmentCheckRequestTupleISCMessage) error {
	if len(tuple) == 0 {
		return nil
	}

	now := time.Now().UTC()
	requests := make([]normalizedKPICheckRequest, 0, len(tuple))
	for _, messagePayload := range tuple {
		eventTime := messagePayload.EventTime
		if eventTime.IsZero() {
			eventTime = now
		}
		params := map[string]interface{}{}
		switch p := messagePayload.Parameters.(type) {
		case map[string]interface{}:
			params = p
		default:
			params = map[string]interface{}{}
		}
		requests = append(requests, normalizedKPICheckRequest{
			Message:   messagePayload,
			Params:    params,
			EventTime: eventTime,
		})
	}

	sort.SliceStable(requests, func(i, j int) bool {
		return requests[i].EventTime.Before(requests[j].EventTime)
	})

	sdTypeUpdates := make(map[string]struct{})
	rawMessages := make([]sharedModel.RawDataPointISCMessage, 0, len(requests))
	rawRecords := make([]sharedModel.TimeSeriesRawRecord, 0, len(requests))
	kpiResults := make([]sharedModel.KPIFulfillmentCheckResultISCMessage, 0, len(requests))
	kpiRecords := make([]sharedModel.TimeSeriesKPIResultRecord, 0, len(requests))

	for _, request := range requests {
		rawOutputs, kpiOutputs, shouldProcessKPI := processRequestWithInstanceLock(request)
		if rawOutputs.SDTypeChanged {
			sdTypeUpdates[request.Message.SDTypeUID] = struct{}{}
		}
		rawMessages = append(rawMessages, rawOutputs.RawMessages...)
		rawRecords = append(rawRecords, rawOutputs.RawRecords...)
		if !shouldProcessKPI {
			continue
		}

		kpiResults = append(kpiResults, kpiOutputs.Results...)
		kpiRecords = append(kpiRecords, kpiOutputs.TimeSeriesKPIRecord...)
	}

	if len(sdTypeUpdates) > 0 {
		sdTypeUIDs := make([]string, 0, len(sdTypeUpdates))
		for sdTypeUID := range sdTypeUpdates {
			sdTypeUIDs = append(sdTypeUIDs, sdTypeUID)
		}
		publishSDTypeUpdates(rabbitMQClient, sdTypeUIDs)
	}
	publishRaw(rabbitMQClient, rawMessages)
	storeRaw(rabbitMQClient, rawRecords)
	storeKPI(rabbitMQClient, kpiRecords)
	publishKPI(rabbitMQClient, kpiResults, false)
	return nil
}

func processRequestWithInstanceLock(request normalizedKPICheckRequest) (RawProcessingOutputs, KPIProcessingOutputs, bool) {
	unlock := lockInstanceForInputProcessing(request.Message.SDInstanceUID)
	defer unlock()

	rawOutputs, shouldProcessKPI := ProcessRaw(request.Message, request.Params, request.EventTime)
	if !shouldProcessKPI {
		return rawOutputs, KPIProcessingOutputs{}, false
	}

	return rawOutputs, ProcessKPI(request.Message, request.Params, request.EventTime), true
}
