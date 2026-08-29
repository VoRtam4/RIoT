/**
 * @file main.go
 * @brief Zpracování interních zpráv Backend Core přijímaných přes RabbitMQ.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ zpracování interní komunikace systému.
 * - Vojtěch Hubáček: doplnění zpracování raw dat, registrace instancí, reprocessingu, mazání, cache zpráv, rozšíření KPI zpráv o vlastní model a tuple zpracování.
 *
 * @ingroup riot_backend_core
 */
package isc

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

var activeReprocessJobs sync.Map

func EnsureIngestQueueForSDType(rabbitMQClient rabbitmq.Client, sdTypeUID string) error {
	if sdTypeUID == "" {
		return fmt.Errorf("cannot declare ingest queue for empty SD type UID")
	}
	queueName := sharedConstants.IngestQueueName(sdTypeUID)
	if err := rabbitMQClient.DeclareExchange(sharedConstants.IngestExchangeName, "direct"); err != nil {
		return err
	}
	if err := rabbitMQClient.DeclareQueue(queueName); err != nil {
		return err
	}
	return rabbitMQClient.BindQueue(queueName, sharedConstants.IngestExchangeName, sdTypeUID)
}

func EnsureIngestQueuesForAllSDTypes(rabbitMQClient rabbitmq.Client) error {
	sdTypesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypes()
	if sdTypesLoadResult.IsFailure() {
		return sdTypesLoadResult.GetError()
	}
	for _, sdType := range sdTypesLoadResult.GetPayload() {
		if err := EnsureIngestQueueForSDType(rabbitMQClient, sdType.UID); err != nil {
			return err
		}
	}
	return nil
}

func ProcessIncomingMessageProcessingUnitConnectionNotifications() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeMessageProcessingUnitConnectionNotificationJSONMessages(func(_ sharedModel.MessageProcessingUnitConnectionNotification) error {
		EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(rabbitMQClient, "")
		EnqueueMessageRepresentingCurrentRawDataCache(rabbitMQClient)
		EnqueueMessageRepresentingCurrentKPICache(rabbitMQClient)
		return nil
	}, rabbitMQClient)
}

func ProcessIncomingSDInstanceRegistrationRequests() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeSDInstanceRegistrationRequestJSONMessages(func(tuple sharedModel.SDInstanceRegistrationRequestTupleISCMessage) error {
		if len(tuple) == 0 {
			return nil
		}
		updated := false
		for _, msg := range tuple {
			log.Printf("SD instance registrace: %s", msg.SDTypeUID)
			sdTypeUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(msg.SDTypeUID, "sdt", "SD type")
			if normalizeErr != nil {
				return normalizeErr
			}
			sdInstanceUID, _, normalizeErr := sharedUtils.NormalizeScopedUID(msg.SDInstanceUID, sdTypeUID, "sdi", "SD instance")
			if normalizeErr != nil {
				return normalizeErr
			}
			result := dbClient.GetRelationalDatabaseClientInstance().UpsertSDInstance(sdInstanceUID, sdTypeUID, msg.Label)
			if result.IsSuccess() {
				instance := result.GetPayload()
				events.GetEventBus().Publish(events.SDInstanceRegisteredEventType, dll2gql.ToGraphQLModelSDInstance(instance))
				updated = true
				continue
			}
			err := result.GetError()
			if errors.Is(err, dbClient.ErrOperationWouldLeadToForeignKeyIntegrityBreach) {
				return fmt.Errorf("SDType not ready yet for uid=%s type=%s", msg.SDInstanceUID, msg.SDTypeUID)
			}
			return fmt.Errorf("failed to persist SD instance (uid=%s type=%s): %w", msg.SDInstanceUID, msg.SDTypeUID, err)
		}
		if updated {
			EnqueueMessageRepresentingCurrentSDInstanceConfiguration(rabbitMQClient)
		}
		return nil
	}, rabbitMQClient)
}

func ProcessIncomingRawDataPoints() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeRawDataPointJSONMessages(func(rawTuple []sharedModel.RawDataPointISCMessage) error {
		if len(rawTuple) == 0 {
			return nil
		}
		grouped := make(map[string][]sharedModel.RawDataPointISCMessage)
		for _, msg := range rawTuple {
			sdTypeUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(msg.SDTypeUID, "sdt", "SD type")
			if normalizeErr != nil {
				return normalizeErr
			}
			sdInstanceUID, _, normalizeErr := sharedUtils.NormalizeScopedUID(msg.SDInstanceUID, sdTypeUID, "sdi", "SD instance")
			if normalizeErr != nil {
				return normalizeErr
			}
			msg.SDTypeUID = sdTypeUID
			msg.SDInstanceUID = sdInstanceUID
			grouped[sdInstanceUID] = append(grouped[sdInstanceUID], msg)
		}
		for uid, messages := range grouped {
			instanceResult := dbClient.GetRelationalDatabaseClientInstance().UpsertSDInstance(uid, messages[0].SDTypeUID, "")
			if instanceResult.IsFailure() {
				return instanceResult.GetError()
			}
			instance := instanceResult.GetPayload()
			instanceID := instance.ID.GetPayload()
			instanceUID := instance.UID
			sdTypeID := instance.SDType.ID.GetPayload()
			sdTypeUID := instance.SDType.UID
			points := make([]dllModel.RawDataPoint, 0, len(messages))
			for _, msg := range messages {
				eventTime := msg.EventTime
				if eventTime.IsZero() {
					eventTime = time.Now().UTC()
				}
				points = append(points, dllModel.RawDataPoint{
					SDTypeID:      sdTypeID,
					SDTypeUID:     sdTypeUID,
					SDInstanceID:  instanceID,
					SDInstanceUID: instanceUID,
					EventTime:     eventTime,
					Payload:       msg.Payload,
				})
			}
			result := dbClient.GetRelationalDatabaseClientInstance().PersistRawDataPoints(points)
			if result.IsFailure() {
				return result.GetError()
			}
			gql := sharedUtils.Map(result.GetPayload(), func(p dllModel.RawDataPoint) graphQLModel.RawDataPoint {
				return graphQLModel.RawDataPoint{
					SdTypeUID:     p.SDTypeUID,
					SdInstanceUID: p.SDInstanceUID,
					EventTime:     p.EventTime.Format(time.RFC3339Nano),
					Payload:       string(p.Payload),
				}
			})
			events.GetEventBus().Publish(events.RawDataPointReceivedEventType, gql)
		}
		return nil
	}, rabbitMQClient)
}

func ProcessIncomingKPIFulfillmentCheckResults() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeKPIFulfillmentCheckResultJSONMessages(func(results sharedModel.KPIFulfillmentCheckResultTupleISCMessage) error {
		if len(results.Tuple) == 0 {
			return nil
		}
		kpiDefinitionsResult := dbClient.GetRelationalDatabaseClientInstance().LoadAllKPIDefinitions()
		if kpiDefinitionsResult.IsFailure() {
			return kpiDefinitionsResult.GetError()
		}
		kpiDefinitionsByUID := make(map[string]sharedModel.KPIDefinition, len(kpiDefinitionsResult.GetPayload()))
		for _, kpiDefinition := range kpiDefinitionsResult.GetPayload() {
			if kpiDefinition.UID != nil {
				kpiDefinitionsByUID[*kpiDefinition.UID] = kpiDefinition
			}
		}
		grouped := make(map[string][]sharedModel.KPIFulfillmentCheckResultISCMessage)
		for _, msg := range results.Tuple {
			sdTypeUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(msg.SDTypeUID, "sdt", "SD type")
			if normalizeErr != nil {
				return normalizeErr
			}
			sdInstanceUID, _, normalizeErr := sharedUtils.NormalizeScopedUID(msg.SDInstanceUID, sdTypeUID, "sdi", "SD instance")
			if normalizeErr != nil {
				return normalizeErr
			}
			msg.SDTypeUID = sdTypeUID
			msg.SDInstanceUID = sdInstanceUID
			grouped[sdInstanceUID] = append(grouped[sdInstanceUID], msg)
		}
		for uid, messages := range grouped {
			instanceResult := dbClient.GetRelationalDatabaseClientInstance().UpsertSDInstance(uid, messages[0].SDTypeUID, "")
			if instanceResult.IsFailure() {
				return instanceResult.GetError()
			}
			instance := instanceResult.GetPayload()
			instanceID := instance.ID.GetPayload()
			instanceUID := instance.UID
			sdTypeID := instance.SDType.ID.GetPayload()
			sdTypeUID := instance.SDType.UID
			points := make([]dllModel.KPIFulfillmentCheckResult, 0, len(messages))
			for _, msg := range messages {
				kpiDefinition, ok := kpiDefinitionsByUID[msg.KPIDefinitionUID]
				if !ok {
					return fmt.Errorf("kpiDefinition not found for UID: %s", msg.KPIDefinitionUID)
				}
				if kpiDefinition.ID == nil {
					return fmt.Errorf("kpiDefinition loaded by UID has no internal ID: %s", msg.KPIDefinitionUID)
				}
				eventTime := msg.EventTime
				if eventTime.IsZero() {
					eventTime = time.Now().UTC()
				}
				points = append(points, dllModel.KPIFulfillmentCheckResult{
					SDTypeID:         sdTypeID,
					SDTypeUID:        sdTypeUID,
					SDInstanceID:     instanceID,
					SDInstanceUID:    instanceUID,
					KPIDefinitionID:  *kpiDefinition.ID,
					KPIDefinitionUID: msg.KPIDefinitionUID,
					Fulfilled:        msg.Fulfilled,
					EventTime:        eventTime,
				})
			}
			result := dbClient.GetRelationalDatabaseClientInstance().PersistKPIFulfillmentCheckResults(points, results.Reprocess)
			if result.IsFailure() {
				return result.GetError()
			}
			gql := sharedUtils.Map(result.GetPayload(), func(p dllModel.KPIFulfillmentCheckResult) graphQLModel.KPIFulfillmentCheckResult {
				return graphQLModel.KPIFulfillmentCheckResult{
					SdTypeUID:        p.SDTypeUID,
					SdInstanceUID:    p.SDInstanceUID,
					KpiDefinitionUID: p.KPIDefinitionUID,
					Fulfilled:        p.Fulfilled,
					EventTime:        p.EventTime.Format(time.RFC3339Nano),
				}
			})
			events.GetEventBus().Publish(events.KPIFulfillmentCheckedEventType, gql)
		}
		return nil
	}, rabbitMQClient)
}

func ProcessIncomingSDTypeRegistrationRequests() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeSDTypeRegistrationRequestJSONMessages(func(tuple sharedModel.SDTypeRegistrationRequestTupleISCMessage) error {
		if len(tuple) == 0 {
			return nil
		}
		for _, message := range tuple {
			sdTypeUID, _, normalizeErr := sharedUtils.NormalizePrefixedUID(message.SDTypeUID, "sdt", "SD type")
			if normalizeErr != nil {
				return normalizeErr
			}
			params := make([]dllModel.SDParameter, 0)
			for _, p := range message.Parameters {
				params = append(params, dllModel.SDParameter{
					Denotation: p.Denotation,
					Label:      sharedUtils.SafeLabel(p.Label, p.Denotation),
					Type:       dllModel.SDParameterType(p.Type),
					Role:       dllModel.SDParameterRole(p.Role),
				})
			}
			sdType := dllModel.SDType{
				UID:        sdTypeUID,
				Label:      sharedUtils.SafeLabel(message.Label, sdTypeUID),
				Parameters: params,
			}
			result := dbClient.GetRelationalDatabaseClientInstance().UpsertSDType(sdType)
			if result.IsFailure() {
				return result.GetError()
			}
			if err := EnsureIngestQueueForSDType(rabbitMQClient, sdTypeUID); err != nil {
				return fmt.Errorf("failed to ensure ingest queue for SDType %s: %w", sdTypeUID, err)
			}
			log.Printf("SDType %s registrován (params=%d).", sdTypeUID, len(params))
		}
		EnqueueMessageRepresentingCurrentSDTypeConfiguration(rabbitMQClient)
		return nil
	}, rabbitMQClient)
}

func EnqueueMessageRepresentingCurrentSDTypeConfiguration(rabbitMQClient rabbitmq.Client) {
	sharedUtils.TerminateOnError(func() error {
		sdTypesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypes()
		if sdTypesLoadResult.IsFailure() {
			return sdTypesLoadResult.GetError()
		}
		sdTypes := sdTypesLoadResult.GetPayload()
		messages := make([]sharedModel.SDTypeUpdateISCMessage, 0)
		for _, sdType := range sdTypes {
			params := make([]sharedModel.SDParameter, 0)
			for _, p := range sdType.Parameters {
				params = append(params, sharedModel.SDParameter{
					Denotation: p.Denotation,
					Label:      sharedUtils.SafeLabel(p.Label, p.Denotation),
					Type:       sharedModel.SDParameterType(p.Type),
					Role:       sharedModel.SDParameterRole(p.Role),
				})
			}
			messages = append(messages, sharedModel.SDTypeUpdateISCMessage{
				SDTypeUID:  sdType.UID,
				Parameters: params,
			})
		}
		jsonResult := sharedUtils.SerializeToJSON(messages)
		if jsonResult.IsFailure() {
			return jsonResult.GetError()
		}
		return rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.SetOfSDTypesUpdatesQueueName), jsonResult.GetPayload())
	}(), "[ISC] Failed to enqueue SDType configuration")
}

func EnqueueMessageRepresentingCurrentSDInstanceConfiguration(rabbitMQClient rabbitmq.Client) {
	sharedUtils.TerminateOnError(func() error {
		sdInstancesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstances()
		if sdInstancesLoadResult.IsFailure() {
			return sdInstancesLoadResult.GetError()
		}
		sdInstancesInfo := sharedUtils.Map[dllModel.SDInstance, sharedModel.SDInstanceInfo](sdInstancesLoadResult.GetPayload(), func(sdInstance dllModel.SDInstance) sharedModel.SDInstanceInfo {
			return sharedModel.SDInstanceInfo{
				SDInstanceUID:   sdInstance.UID,
				ConfirmedByUser: sdInstance.ConfirmedByUser,
			}
		})
		sdInstancesInfoJSONSerializationResult := sharedUtils.SerializeToJSON(sdInstancesInfo)
		if sdInstancesInfoJSONSerializationResult.IsFailure() {
			return sdInstancesInfoJSONSerializationResult.GetError()
		}
		return rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.SetOfSDInstancesUpdatesQueueName), sdInstancesInfoJSONSerializationResult.GetPayload())
	}(), "[ISC] Failed to enqueue RabbitMQ messages representing current SD instance configuration")
}

func EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(rabbitMQClient rabbitmq.Client, jobID string) {
	sharedUtils.TerminateOnError(func() error {
		kpiDefinitionsLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadAllKPIDefinitions()
		if kpiDefinitionsLoadResult.IsFailure() {
			return kpiDefinitionsLoadResult.GetError()
		}
		kpiDefinitions := kpiDefinitionsLoadResult.GetPayload()
		kpiDefinitionsBySDTypeDenotationMap := make(map[string][]sharedModel.KPIDefinitionMPU)
		for _, kpiDefinition := range kpiDefinitions {
			kpiDefinitionUID := sharedUtils.NewOptionalFromPointer(kpiDefinition.UID).GetPayload()
			if kpiDefinitionUID == "" {
				continue
			}
			sdTypeSpecification := kpiDefinition.SDTypeSpecification
			if _, exists := kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification]; !exists {
				kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification] = make([]sharedModel.KPIDefinitionMPU, 0)
			}
			var selectedSDInstanceUIDs []string
			for _, sdInstanceID := range kpiDefinition.SelectedSDInstanceIDs {
				uid := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstance(sdInstanceID)
				if uid.IsFailure() {
					continue
				}
				selectedSDInstanceUIDs = append(selectedSDInstanceUIDs, uid.GetPayload().UID)
			}
			kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification] = append(kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification], sharedModel.KPIDefinitionMPU{
				UID:                    kpiDefinitionUID,
				SDTypeUID:              kpiDefinition.SDTypeSpecification,
				RootNode:               kpiDefinition.RootNode,
				SDInstanceMode:         kpiDefinition.SDInstanceMode,
				SelectedSDInstanceUIDs: selectedSDInstanceUIDs,
			})
		}
		message := sharedModel.KPIConfigurationUpdateISCMessage{
			KpiConfiguration: kpiDefinitionsBySDTypeDenotationMap,
			JobID:            jobID,
		}
		jsonSerializationResult := sharedUtils.SerializeToJSON(message)
		if jsonSerializationResult.IsFailure() {
			return jsonSerializationResult.GetError()
		}
		return rabbitMQClient.PublishJSONMessage(sharedUtils.NewOptionalOf(sharedConstants.BuiltInFanoutExchangeName), sharedUtils.NewEmptyOptional[string](), jsonSerializationResult.GetPayload())
	}(), "[ISC] Failed to enqueue RabbitMQ messages representing current KPI definition configuration")
}

func EnqueueMessageRepresentingCurrentRawDataCache(rabbitMQClient rabbitmq.Client) {
	const batchSize = 500
	sharedUtils.TerminateOnError(func() error {
		rawPointsLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadAllRawDataPoints()
		if rawPointsLoadResult.IsFailure() {
			return rawPointsLoadResult.GetError()
		}
		sdInstancesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstances()
		if sdInstancesLoadResult.IsFailure() {
			return sdInstancesLoadResult.GetError()
		}

		instanceByID := make(map[uint32]dllModel.SDInstance, len(sdInstancesLoadResult.GetPayload()))
		for _, instance := range sdInstancesLoadResult.GetPayload() {
			id := instance.ID.GetPayloadOrDefault(0)
			if id == 0 {
				continue
			}
			instanceByID[id] = instance
		}

		messages := make([]sharedModel.RawDataPointISCMessage, 0, len(rawPointsLoadResult.GetPayload()))
		for _, point := range rawPointsLoadResult.GetPayload() {
			instance, exists := instanceByID[point.SDInstanceID]
			if !exists {
				continue
			}
			messages = append(messages, sharedModel.RawDataPointISCMessage{
				SDTypeUID:     instance.SDType.UID,
				SDInstanceUID: instance.UID,
				EventTime:     point.EventTime,
				Payload:       point.Payload,
			})
		}

		if len(messages) == 0 {
			payloadResult := sharedUtils.SerializeToJSON(sharedModel.RawDataPointCacheBootstrapISCMessage{
				Tuple: []sharedModel.RawDataPointISCMessage{},
				Done:  true,
			})
			if payloadResult.IsFailure() {
				return payloadResult.GetError()
			}
			return rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.RawDataPointCacheBootstrapQueueName), payloadResult.GetPayload())
		}

		for start := 0; start < len(messages); start += batchSize {
			end := start + batchSize
			if end > len(messages) {
				end = len(messages)
			}
			payloadResult := sharedUtils.SerializeToJSON(sharedModel.RawDataPointCacheBootstrapISCMessage{
				Tuple: messages[start:end],
				Done:  end == len(messages),
			})
			if payloadResult.IsFailure() {
				return payloadResult.GetError()
			}
			if err := rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.RawDataPointCacheBootstrapQueueName), payloadResult.GetPayload()); err != nil {
				return err
			}
		}

		return nil
	}(), "[ISC] Failed to enqueue RAW cache bootstrap")
}

func EnqueueMessageRepresentingCurrentKPICache(rabbitMQClient rabbitmq.Client) {
	const batchSize = 500
	sharedUtils.TerminateOnError(func() error {
		kpiResultsLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadAllKPIFulfillmentCheckResults()
		if kpiResultsLoadResult.IsFailure() {
			return kpiResultsLoadResult.GetError()
		}
		sdInstancesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstances()
		if sdInstancesLoadResult.IsFailure() {
			return sdInstancesLoadResult.GetError()
		}

		instanceByID := make(map[uint32]dllModel.SDInstance, len(sdInstancesLoadResult.GetPayload()))
		for _, instance := range sdInstancesLoadResult.GetPayload() {
			id := instance.ID.GetPayloadOrDefault(0)
			if id == 0 {
				continue
			}
			instanceByID[id] = instance
		}

		messages := make([]sharedModel.KPIFulfillmentCheckResultISCMessage, 0, len(kpiResultsLoadResult.GetPayload()))
		for _, result := range kpiResultsLoadResult.GetPayload() {
			instance, exists := instanceByID[result.SDInstanceID]
			if !exists {
				continue
			}
			messages = append(messages, sharedModel.KPIFulfillmentCheckResultISCMessage{
				SDTypeUID:        instance.SDType.UID,
				SDInstanceUID:    instance.UID,
				KPIDefinitionUID: result.KPIDefinitionUID,
				EventTime:        result.EventTime,
				Fulfilled:        result.Fulfilled,
			})
		}

		if len(messages) == 0 {
			payloadResult := sharedUtils.SerializeToJSON(sharedModel.KPIFulfillmentCacheBootstrapISCMessage{
				Tuple: []sharedModel.KPIFulfillmentCheckResultISCMessage{},
				Done:  true,
			})
			if payloadResult.IsFailure() {
				return payloadResult.GetError()
			}
			return rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.KPIFulfillmentCacheBootstrapQueueName), payloadResult.GetPayload())
		}

		for start := 0; start < len(messages); start += batchSize {
			end := start + batchSize
			if end > len(messages) {
				end = len(messages)
			}
			payloadResult := sharedUtils.SerializeToJSON(sharedModel.KPIFulfillmentCacheBootstrapISCMessage{
				Tuple: messages[start:end],
				Done:  end == len(messages),
			})
			if payloadResult.IsFailure() {
				return payloadResult.GetError()
			}
			if err := rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.KPIFulfillmentCacheBootstrapQueueName), payloadResult.GetPayload()); err != nil {
				return err
			}
		}

		return nil
	}(), "[ISC] Failed to enqueue KPI cache bootstrap")
}

func EnqueueMessagesRepresentingCurrentSystemConfiguration(rabbitMQClient rabbitmq.Client) {
	EnqueueMessageRepresentingCurrentSDTypeConfiguration(rabbitMQClient)
	EnqueueMessageRepresentingCurrentSDInstanceConfiguration(rabbitMQClient)
	EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(rabbitMQClient, "")
}

func EnqueueKPIReprocessRequest(rabbitMQClient rabbitmq.Client, kpiDefinition sharedModel.KPIDefinition, fromTime time.Time, jobID string, wait bool) error {
	req := sharedModel.KPIReprocessRequestISCMessage{
		JobID:            jobID,
		Wait:             wait,
		KPIDefinitionUID: sharedUtils.NewOptionalFromPointer(kpiDefinition.UID).GetPayload(),
		SDTypeUID:        kpiDefinition.SDTypeSpecification,
		To:               fromTime,
	}
	if kpiDefinition.SDInstanceMode == sharedModel.SELECTED {
		var selectedSDInstanceUIDs []string
		for _, sdInstanceID := range kpiDefinition.SelectedSDInstanceIDs {
			uid := dbClient.GetRelationalDatabaseClientInstance().LoadSDInstance(sdInstanceID)
			if uid.IsFailure() {
				continue
			}
			selectedSDInstanceUIDs = append(selectedSDInstanceUIDs, uid.GetPayload().UID)
		}
		req.SDInstanceUIDs = selectedSDInstanceUIDs
	}
	jsonResult := sharedUtils.SerializeToJSON(req)
	if jsonResult.IsFailure() {
		return jsonResult.GetError()
	}
	return rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.KPIReprocessRequestQueueName), jsonResult.GetPayload())
}

func EnqueueKPIDeleteRequest(rabbitMQClient rabbitmq.Client, kpiDefinition sharedModel.KPIDefinition, jobID string) error {
	payload := sharedModel.KPIDeleteResultsRequestISCMessage{
		JobID:            jobID,
		KPIDefinitionUID: sharedUtils.NewOptionalFromPointer(kpiDefinition.UID).GetPayload(),
		SDTypeUID:        kpiDefinition.SDTypeSpecification,
	}
	jsonResult := sharedUtils.SerializeToJSON(payload)
	if jsonResult.IsFailure() {
		return jsonResult.GetError()
	}
	if err := rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TSDBDeleteQueueName), jsonResult.GetPayload()); err != nil {
		return err
	}
	if err := rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.MPUDeleteQueueName), jsonResult.GetPayload()); err != nil {
		return err
	}
	return nil
}
