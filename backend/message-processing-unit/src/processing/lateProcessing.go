package processing

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessLateRecord(client rabbitmq.Client, messagePayload sharedModel.KPIFulfillmentCheckRequestISCMessage, params map[string]interface{}, eventTime time.Time) error {
	neighborhood, err := getRecordNeighborhood(client, sharedModel.TimeSeriesRecordNeighborhoodRequest{
		SDTypeUID:     messagePayload.SDTypeUID,
		SDInstanceUID: messagePayload.SDInstanceUID,
		EventTime:     eventTime,
	})
	if err != nil {
		return err
	}
	previousValues := map[string]interface{}{}
	if neighborhood.Previous != nil {
		previousValues = mergeTagsAndFields(*neighborhood.Previous)
	}
	incomingValues := processParams(messagePayload.SDTypeUID, params, previousValues).Values
	if neighborhood.SameTimestamp != nil {
		sameValues := mergeTagsAndFields(*neighborhood.SameTimestamp)
		if sharedUtils.CompareJSONs(sameValues, incomingValues) {
			return nil
		}
	}

	currentRecords, currentState := evaluateKPIRecords(messagePayload.SDTypeUID, messagePayload.SDInstanceUID, incomingValues, previousValues, eventTime)

	currentOperation := sharedModel.TimeSeriesRecordOperationUpsert
	if neighborhood.SameTimestamp == nil && neighborhood.Previous != nil && statesEquivalent(incomingValues, currentState, previousValues, neighborhood.PreviousKPIState) {
		currentOperation = sharedModel.TimeSeriesRecordOperationNone
	}

	correction := sharedModel.TimeSeriesLateRecordCorrection{
		CorrectionID:       uuid.New().String(),
		SDTypeUID:          messagePayload.SDTypeUID,
		SDInstanceUID:      messagePayload.SDInstanceUID,
		EventTime:          eventTime,
		CurrentOperation:   currentOperation,
		SuccessorOperation: sharedModel.TimeSeriesRecordOperationNone,
		CurrentKPIRecords:  currentRecords,
	}
	if currentOperation == sharedModel.TimeSeriesRecordOperationUpsert {
		fields, tags := splitParamsBySDType(messagePayload.SDTypeUID, incomingValues)
		correction.CurrentRawRecord = &sharedModel.TimeSeriesRawRecord{
			EventTime:     eventTime,
			SDInstanceUID: messagePayload.SDInstanceUID,
			SDTypeUID:     messagePayload.SDTypeUID,
			Fields:        fields,
			Tags:          tags,
		}
	}

	if neighborhood.Next != nil {
		nextValues := mergeTagsAndFields(*neighborhood.Next)
		successorRecords, successorState := evaluateKPIRecords(messagePayload.SDTypeUID, messagePayload.SDInstanceUID, nextValues, incomingValues, neighborhood.Next.Time)
		successorTime := neighborhood.Next.Time
		correction.SuccessorRawTime = &successorTime
		if statesEquivalent(nextValues, successorState, incomingValues, currentState) {
			correction.SuccessorOperation = sharedModel.TimeSeriesRecordOperationDelete
		} else {
			correction.SuccessorOperation = sharedModel.TimeSeriesRecordOperationUpsert
			correction.SuccessorKPIRecords = successorRecords
		}
	}

	return publishLateRecordCorrection(client, correction)
}

func getRecordNeighborhood(client rabbitmq.Client, req sharedModel.TimeSeriesRecordNeighborhoodRequest) (sharedModel.TimeSeriesRecordNeighborhoodResponse, error) {
	ch := client.GetChannel()
	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, err
	}
	msgs, err := ch.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, err
	}
	correlationID := uuid.New().String()
	body, err := json.Marshal(req)
	if err != nil {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := ch.PublishWithContext(ctx, "", sharedConstants.TimeSeriesRecordNeighborhoodRequestQueueName, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Body:          body,
		CorrelationId: correlationID,
		ReplyTo:       replyQueue.Name,
	}); err != nil {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, err
	}

	var response sharedModel.TimeSeriesRecordNeighborhoodResponse
	err = rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesRecordNeighborhoodResponse](msgs, correlationID, 10*time.Second, func(resp sharedModel.TimeSeriesRecordNeighborhoodResponse, msg amqp.Delivery) (bool, error) {
		response = resp
		return true, nil
	})
	if err != nil {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, err
	}
	if response.Error != "" {
		return sharedModel.TimeSeriesRecordNeighborhoodResponse{}, errors.New(response.Error)
	}
	return response, nil
}

func publishLateRecordCorrection(client rabbitmq.Client, correction sharedModel.TimeSeriesLateRecordCorrection) error {
	payloadResult := sharedUtils.SerializeToJSON(correction)
	if payloadResult.IsFailure() {
		return payloadResult.GetError()
	}
	return client.PublishJSONMessage(
		sharedUtils.NewEmptyOptional[string](),
		sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesLateRecordCorrectionQueueName),
		payloadResult.GetPayload(),
	)
}

func evaluateKPIRecords(sdTypeUID string, sdInstanceUID string, currentRaw map[string]interface{}, previousRaw map[string]interface{}, eventTime time.Time) ([]sharedModel.TimeSeriesKPIResultRecord, map[string]bool) {
	kpiDefinitionsBySDTypeDenotationMapMutex.RLock()
	kpiDefinitions := kpiDefinitionsBySDTypeDenotationMap[sdTypeUID]
	kpiDefinitionsBySDTypeDenotationMapMutex.RUnlock()
	_, tags := splitParamsBySDType(sdTypeUID, currentRaw)
	context := KPIEvaluationContext{
		CurrentValues:  currentRaw,
		PreviousValues: previousRaw,
	}
	records := make([]sharedModel.TimeSeriesKPIResultRecord, 0, len(kpiDefinitions))
	state := make(map[string]bool, len(kpiDefinitions))
	for _, kpiDefinition := range kpiDefinitions {
		if !kpiAppliesToInstance(kpiDefinition, sdInstanceUID) {
			continue
		}
		result := CheckKPIFulfillment(kpiDefinition, context)
		if result.IsFailure() {
			continue
		}
		value := result.GetPayload()
		state[kpiDefinition.UID] = value
		records = append(records, sharedModel.TimeSeriesKPIResultRecord{
			EventTime:        eventTime,
			SDInstanceUID:    sdInstanceUID,
			SDTypeUID:        sdTypeUID,
			KPIDefinitionUID: kpiDefinition.UID,
			Fulfilled:        value,
			Tags:             tags,
		})
	}
	return records, state
}

func kpiAppliesToInstance(kpiDefinition sharedModel.KPIDefinitionMPU, sdInstanceUID string) bool {
	if kpiDefinition.SDInstanceMode == sharedModel.ALL {
		return true
	}
	for _, uid := range kpiDefinition.SelectedSDInstanceUIDs {
		if uid == sdInstanceUID {
			return true
		}
	}
	return false
}

func statesEquivalent(leftRaw map[string]interface{}, leftKPI map[string]bool, rightRaw map[string]interface{}, rightKPI map[string]bool) bool {
	return sharedUtils.CompareJSONs(leftRaw, rightRaw) && sharedUtils.CompareJSONs(leftKPI, rightKPI)
}
