/**
 * @file setup.go
 * @brief Příprava RabbitMQ infrastruktury používané interní komunikací Backend Core.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ deklarace RabbitMQ infrastruktury.
 * - Vojtěch Hubáček: doplnění deklarací front vycházejících z nových sdílených konstant pro historizaci, cache, reprocessing a registraci instancí.
 *
 * @ingroup riot_backend_core
 */
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
		sharedConstants.RawDataPointCacheBootstrapQueueName,
		sharedConstants.KPIFulfillmentCacheBootstrapQueueName,
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
		sharedConstants.TimeSeriesReadCancelRequestQueueName,
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
