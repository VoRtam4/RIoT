package isc

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func getQueueDeclarationErrorMessage(queueName string) string {
	return fmt.Sprintf("[ISC setup] Failed to declare the '%s' queue", queueName)
}

func SetupRabbitMQInfrastructureForISC(rabbitMQClient rabbitmq.Client) {
	namesOfQueuesToDeclare := sharedUtils.SliceOf[string](
		sharedConstants.KPIFulfillmentCheckResultsQueueName,
		sharedConstants.RawDataPointQueueName,
		sharedConstants.KPIFulfillmentCheckRequestsQueueName,
		sharedConstants.SDInstanceRegistrationRequestsQueueName,
		sharedConstants.SDTypeRegistrationRequestsQueueName,
		sharedConstants.SetOfSDTypesUpdatesQueueName,
		sharedConstants.SetOfSDInstancesUpdatesQueueName,
		sharedConstants.MessageProcessingUnitConnectionNotificationsQueueName,
		sharedConstants.TimeSeriesRawDataQueueName,
		sharedConstants.TimeSeriesKPIResultQueueName,
		sharedConstants.TimeSeriesReadRequestQueueName,
		sharedConstants.TimeSeriesReadResponseQueueName,
		sharedConstants.TimeSeriesDistinctTagValuesRequestQueueName,
		sharedConstants.TimeSeriesDistinctTagValuesResponseQueueName,
		sharedConstants.TimeSeriesReprocessReadRequestQueueName,
		sharedConstants.TimeSeriesReprocessReadResponseQueueName,
		sharedConstants.KPIReprocessRequestQueueName,
		sharedConstants.TSDBDeleteQueueName,
		sharedConstants.MPUDeleteQueueName,
		sharedConstants.KPIConfigUpdateQueueName,
	)
	sharedUtils.ForEach(namesOfQueuesToDeclare, func(nameOfQueuesToDeclare string) {
		sharedUtils.TerminateOnError(rabbitMQClient.DeclareQueue(nameOfQueuesToDeclare), getQueueDeclarationErrorMessage(nameOfQueuesToDeclare))
	})
}
