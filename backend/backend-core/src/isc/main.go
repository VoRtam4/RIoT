package isc

import (
	"errors"
	"fmt"
	"log"
	"strings"
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

func ProcessIncomingMessageProcessingUnitConnectionNotifications() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeMessageProcessingUnitConnectionNotificationJSONMessages(func(_ sharedModel.MessageProcessingUnitConnectionNotification) error {
		EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(rabbitMQClient)
		return nil
	}, rabbitMQClient)
}

func ProcessIncomingSDInstanceRegistrationRequests() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeSDInstanceRegistrationRequestJSONMessages(func(msg sharedModel.SDInstanceRegistrationRequestISCMessage) error {
		uid := msg.SDInstanceUID
		sdType := msg.SDTypeSpecification
		log.Printf("[SDINSTANCE] Processing registration | uid=%s type=%s", uid, sdType)
		result := dbClient.GetRelationalDatabaseClientInstance().PersistNewSDInstance(uid, sdType)
		if result.IsSuccess() {
			instance := result.GetPayload()
			events.GetEventBus().Publish(events.SDInstanceRegisteredEventType, dll2gql.ToGraphQLModelSDInstance(instance))
			log.Printf("[SDINSTANCE] Registered successfully | uid=%s type=%s", uid, sdType)
			return nil
		}
		err := result.GetError()
		if errors.Is(err, dbClient.ErrOperationWouldLeadToForeignKeyIntegrityBreach) {
			log.Printf("[SDINSTANCE] FK constraint skip | uid=%s type=%s", uid, sdType)
			return nil
		}
		return fmt.Errorf("failed to persist SD instance (uid=%s type=%s): %w", uid, sdType, err)
	}, rabbitMQClient)
}

func ProcessIncomingKPIFulfillmentCheckResults() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	consumeKPIFulfillmentCheckResultJSONMessages(func(kpiFulfillmentCheckResultTuple sharedModel.KPIFulfillmentCheckResultTupleISCMessage) error {
		if len(kpiFulfillmentCheckResultTuple) == 0 {
			log.Printf("Warning: Got an enpty tuple of KPI fulfillment check results... that shouldn't happen...")
			return nil
		}
		eventTimes := make([]time.Time, 0)
		sdInstanceUID := kpiFulfillmentCheckResultTuple[0].SDInstanceUID
		kpiDefinitionIDs := make([]uint32, 0)
		fulfillmentStatuses := make([]bool, 0)
		for _, kpiFulfillmentCheckResult := range kpiFulfillmentCheckResultTuple {
			kpiDefinitionIDs = append(kpiDefinitionIDs, kpiFulfillmentCheckResult.KPIDefinitionID)
			fulfillmentStatuses = append(fulfillmentStatuses, kpiFulfillmentCheckResult.Fulfilled)
			eventTime := kpiFulfillmentCheckResult.EventTime
			if eventTime.IsZero() {
				eventTime = time.Now().UTC()
			}
			eventTimes = append(eventTimes, eventTime)
		}
		kpiFulfillmentCheckResultTuplePersistResult := dbClient.GetRelationalDatabaseClientInstance().PersistKPIFulFulfillmentCheckResultTuple(sdInstanceUID, kpiDefinitionIDs, fulfillmentStatuses, eventTimes)
		if kpiFulfillmentCheckResultTuplePersistResult.IsSuccess() {
			gqlKPIFulfillmentCheckResultTuple := graphQLModel.KPIFulfillmentCheckResultTuple{
				KpiFulfillmentCheckResults: sharedUtils.Map(kpiFulfillmentCheckResultTuplePersistResult.GetPayload(), dll2gql.ToGraphQLModelKPIFulfillmentCheckResult),
			}
			events.GetEventBus().Publish(events.KPIFulfillmentCheckedEventType, gqlKPIFulfillmentCheckResultTuple)
		} else if kpiFulfillmentCheckResultTuplePersistError := kpiFulfillmentCheckResultTuplePersistResult.GetError(); !errors.Is(kpiFulfillmentCheckResultTuplePersistError, dbClient.ErrOperationWouldLeadToForeignKeyIntegrityBreach) {
			return fmt.Errorf("failed to persist KPI fulfillment check result tuple: %w", kpiFulfillmentCheckResultTuplePersistError)
		}
		return nil
	}, rabbitMQClient)
}

func ProcessIncomingSDTypeRegistrationRequests() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	rabbitmq.ConsumeJSONMessages[sharedModel.SDTypeRegistrationRequestISCMessage](rabbitMQClient, sharedConstants.SDTypeRegistrationRequestsQueueName,
		func(message sharedModel.SDTypeRegistrationRequestISCMessage) error {
			params := make([]dllModel.SDParameter, 0)
			for _, p := range message.Parameters {
				var paramType dllModel.SDParameterType
				switch strings.ToUpper(p.Type) {
				case "STRING":
					paramType = dllModel.SDParameterTypeString
				case "NUMBER":
					paramType = dllModel.SDParameterTypeNumber
				case "BOOLEAN":
					paramType = dllModel.SDParameterTypeBoolean
				default:
					return fmt.Errorf("unknown parameter type: %s", p.Type)
				}
				var paramRole dllModel.SDParameterRole
				switch strings.ToUpper(p.Role) {
				case "FIELD":
					paramRole = dllModel.SDParameterRoleField
				case "TAG":
					paramRole = dllModel.SDParameterRoleTag
				default:
					return fmt.Errorf("unknown parameter role: %s", p.Role)
				}
				params = append(params, dllModel.SDParameter{
					Denotation: p.Denotation,
					Type:       paramType,
					Role:       paramRole,
				})
			}
			sdType := dllModel.SDType{
				Denotation: message.SDTypeSpecification,
				Parameters: params,
			}
			result := dbClient.GetRelationalDatabaseClientInstance().UpsertSDType(sdType)
			if result.IsFailure() {
				return result.GetError()
			}
			EnqueueMessageRepresentingCurrentSDTypeConfiguration(rabbitMQClient)
			log.Printf("SDType %s registrován (params=%d).", message.SDTypeSpecification, len(params))
			return nil
		},
	)
}

func EnqueueMessageRepresentingCurrentSDTypeConfiguration(rabbitMQClient rabbitmq.Client) {
	sharedUtils.TerminateOnError(func() error {
		sdTypesLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadSDTypes()
		if sdTypesLoadResult.IsFailure() {
			return sdTypesLoadResult.GetError()
		}
		sdTypes := sdTypesLoadResult.GetPayload()
		messages := make([]sharedModel.SDTypeRegistrationRequestISCMessage, 0)
		for _, sdType := range sdTypes {
			params := make([]sharedModel.SDParameter, 0)
			for _, p := range sdType.Parameters {
				params = append(params, sharedModel.SDParameter{
					Denotation: p.Denotation,
					Type:       string(p.Type),
					Role:       string(p.Role),
				})
			}
			messages = append(messages, sharedModel.SDTypeRegistrationRequestISCMessage{
				SDTypeSpecification: sdType.Denotation,
				Parameters:          params,
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

func EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(rabbitMQClient rabbitmq.Client) {
	sharedUtils.TerminateOnError(func() error {
		kpiDefinitionsLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadKPIDefinitions()
		if kpiDefinitionsLoadResult.IsFailure() {
			return kpiDefinitionsLoadResult.GetError()
		}
		kpiDefinitions := kpiDefinitionsLoadResult.GetPayload()
		kpiDefinitionsBySDTypeDenotationMap := make(sharedModel.KPIConfigurationUpdateISCMessage)
		for _, kpiDefinition := range kpiDefinitions {
			sdTypeSpecification := kpiDefinition.SDTypeSpecification
			if _, exists := kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification]; !exists {
				kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification] = make([]sharedModel.KPIDefinition, 0)
			}
			kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification] = append(kpiDefinitionsBySDTypeDenotationMap[sdTypeSpecification], kpiDefinition)
		}
		jsonSerializationResult := sharedUtils.SerializeToJSON(kpiDefinitionsBySDTypeDenotationMap)
		if jsonSerializationResult.IsFailure() {
			return jsonSerializationResult.GetError()
		}
		return rabbitMQClient.PublishJSONMessage(sharedUtils.NewOptionalOf(sharedConstants.BuiltInFanoutExchangeName), sharedUtils.NewEmptyOptional[string](), jsonSerializationResult.GetPayload())
	}(), "[ISC] Failed to enqueue RabbitMQ messages representing current KPI definition configuration")
}

func EnqueueMessagesRepresentingCurrentSystemConfiguration(rabbitMQClient rabbitmq.Client) {
	EnqueueMessageRepresentingCurrentSDTypeConfiguration(rabbitMQClient)
	EnqueueMessageRepresentingCurrentSDInstanceConfiguration(rabbitMQClient)
	EnqueueMessageRepresentingCurrentKPIDefinitionConfiguration(rabbitMQClient)
}

func EnqueueKPIReprocessRequest(rabbitMQClient rabbitmq.Client, kpiDefinition sharedModel.KPIDefinition, fromTime time.Time) error {
	req := sharedModel.KPIReprocessRequestISCMessage{
		KPIDefinitionID:     sharedUtils.NewOptionalFromPointer(kpiDefinition.ID).GetPayload(),
		SDTypeSpecification: kpiDefinition.SDTypeSpecification,
		To:                  fromTime,
	}
	if kpiDefinition.SDInstanceMode == sharedModel.SELECTED {
		req.SDInstanceUIDs = kpiDefinition.SelectedSDInstanceUIDs
	}
	jsonResult := sharedUtils.SerializeToJSON(req)
	if jsonResult.IsFailure() {
		return jsonResult.GetError()
	}
	return rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.KPIReprocessRequestQueueName), jsonResult.GetPayload())
}

func EnqueueKPIDeleteRequest(rabbitMQClient rabbitmq.Client, kpiDefinitionID uint32) error {
	jsonResult := sharedUtils.SerializeToJSON(sharedModel.KPIDeleteResultsRequestISCMessage{
		KPIDefinitionID: kpiDefinitionID,
	})
	if jsonResult.IsFailure() {
		return jsonResult.GetError()
	}
	return rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesDeleteRequestQueueName), jsonResult.GetPayload())
}
