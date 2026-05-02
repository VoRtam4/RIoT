package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/MichalBures-OG/bp-bures-RIoT-message-processing-unit/src/processing"
)

func checkForKPIFulfilmentCheckRequests() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	err := rabbitmq.ConsumeJSONMessages[sharedModel.KPIFulfillmentCheckRequestTupleISCMessage](
		rabbitMQClient,
		sharedConstants.KPIFulfillmentCheckRequestsQueueName,
		func(tuple sharedModel.KPIFulfillmentCheckRequestTupleISCMessage) error {
			return processing.ProcessKPIFulfillmentCheckRequestTuple(rabbitMQClient, tuple)
		},
	)
	if err != nil {
		log.Printf(
			"Consumption failed from '%s': %s",
			sharedConstants.KPIFulfillmentCheckRequestsQueueName,
			err.Error(),
		)
	}
}

func checkForKPIReprocessRequests() {
	for {
		rabbitMQClient := rabbitmq.NewClient()
		err := rabbitmq.ConsumeJSONMessages[sharedModel.KPIReprocessRequestISCMessage](rabbitMQClient, sharedConstants.KPIReprocessRequestQueueName,
			func(req sharedModel.KPIReprocessRequestISCMessage) error {
				log.Printf("Starting KPI reprocess for KPI definition %d", req.KPIDefinitionID)
				return processing.ReprocessKPI(req)
			},
		)
		rabbitMQClient.Dispose()
		if err != nil {
			log.Printf("Failed consuming KPI reprocess queue: %s", err.Error())
		}
		time.Sleep(time.Second)
	}
}

func checkForKPIDefinitionsBySDTypeDenotationMapUpdates() {
	go func() {
		time.Sleep(time.Second)
		rabbitMQClient := rabbitmq.NewClient()
		defer rabbitMQClient.Dispose()
		if err := rabbitMQClient.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.MessageProcessingUnitConnectionNotificationsQueueName), []byte("{}")); err != nil {
			log.Printf("Failed to publish the message processing unit connection notification message to the %s queue: %s\n", sharedConstants.MessageProcessingUnitConnectionNotificationsQueueName, err.Error())
		}
	}()
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	processing.DenotationMapUpdates(rabbitMQClient)
}

func checkForSDTypeUpdates() {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	err := rabbitmq.ConsumeJSONMessages[[]sharedModel.SDTypeUpdateISCMessage](rabbitMQClient, sharedConstants.SetOfSDTypesUpdatesQueueName,
		func(messages []sharedModel.SDTypeUpdateISCMessage) error {
			processing.UpdateSDType(messages)
			return nil
		},
	)
	if err != nil {
		log.Printf("[MPU][SDTYPE] Failed consuming SDType updates: %s", err.Error())
	}
}

func checkForKPIDeleteRequests() {
	for {
		rabbitMQClient := rabbitmq.NewClient()
		err := rabbitmq.ConsumeJSONMessages[sharedModel.KPIDeleteResultsRequestISCMessage](rabbitMQClient, sharedConstants.MPUDeleteQueueName, processing.ProcessDelete)
		rabbitMQClient.Dispose()
		if err != nil {
			log.Printf("[MPU][DELETE] Failed consuming delete requests: %s", err.Error())
		}
		time.Sleep(time.Second)
	}
}

func main() {
	log.SetOutput(os.Stderr)
	log.Println("Waiting for dependencies...")
	rawBackendCoreURL := sharedUtils.GetEnvironmentVariableValue("BACKEND_CORE_URL").GetPayloadOrDefault("http://riot-backend-core:9090")
	parsedBackendCoreURL, err := url.Parse(rawBackendCoreURL)
	sharedUtils.TerminateOnError(err, fmt.Sprintf("Unable to parse the backend-core URL: %s", rawBackendCoreURL))
	sharedUtils.TerminateOnError(sharedUtils.WaitForDSs(time.Minute, sharedUtils.NewPairOf(parsedBackendCoreURL.Hostname(), parsedBackendCoreURL.Port())), "Some dependencies of this application are inaccessible")
	log.Println("Dependencies should be up and running...")
	processing.InitializeProcessing()
	sharedUtils.StartLoggingProfilingInformationPeriodically(time.Minute)
	sharedUtils.WaitForAll(checkForKPIDefinitionsBySDTypeDenotationMapUpdates, checkForKPIFulfilmentCheckRequests, checkForKPIReprocessRequests, checkForSDTypeUpdates, checkForKPIDeleteRequests)
	processing.StartCacheCleanup(25*time.Hour, time.Hour)
}
