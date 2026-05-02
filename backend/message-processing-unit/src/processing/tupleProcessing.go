package processing

import (
	"sort"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

type normalizedKPICheckRequest struct {
	Message   sharedModel.KPIFulfillmentCheckRequestISCMessage
	Params    map[string]interface{}
	EventTime time.Time
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
		rawOutputs, shouldProcessKPI := ProcessRaw(request.Message, request.Params, request.EventTime)
		if rawOutputs.SDTypeChanged {
			sdTypeUpdates[request.Message.SDTypeUID] = struct{}{}
		}
		rawMessages = append(rawMessages, rawOutputs.RawMessages...)
		rawRecords = append(rawRecords, rawOutputs.RawRecords...)
		if !shouldProcessKPI {
			continue
		}

		kpiOutputs := ProcessKPI(request.Message, request.Params, request.EventTime)
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
