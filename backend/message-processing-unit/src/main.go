/**
 * @file main.go
 * @brief Vstupní bod modulu Message Processing Unit platformy RIoT a registrace jeho konzumentů zpráv.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ vstupního bodu a spuštění hlavního toku modulu.
 * - Vojtěch Hubáček: doplnění workerů a posluchačů pro aktualizace typů zdrojů, reprocessing, rušení úloh a navazující provozní logiku.
 *
 * @defgroup riot_message_processing_unit Message Processing Unit
 * @ingroup riot
 * @see ../README.md
 */
package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/MichalBures-OG/bp-bures-RIoT-message-processing-unit/src/dispatch"
	"github.com/MichalBures-OG/bp-bures-RIoT-message-processing-unit/src/processing"
)

func getWorkerCount(envName string, logPrefix string) int {
	raw := sharedUtils.GetEnvironmentVariableValue(envName).GetPayloadOrDefault("1")
	workerCount, convertErr := strconv.Atoi(raw)
	if convertErr != nil || workerCount < 1 {
		log.Printf("%s Invalid %s=%q, using 1", logPrefix, envName, raw)
		return 1
	}
	return workerCount
}

func checkForKPIFulfilmentCheckRequests() {
	consumerCount := getWorkerCount("MPU_INPUT_WORKERS", "[MPU][INPUT]")
	log.Printf("[MPU][INPUT] Starting %d legacy input consumer(s) | queue=%s", consumerCount, sharedConstants.KPIFulfillmentCheckRequestsQueueName)
	var wg sync.WaitGroup
	for consumerID := 1; consumerID <= consumerCount; consumerID++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for {
				rabbitMQClient := rabbitmq.NewClient()
				err := dispatch.ConsumeInputQueue(rabbitMQClient, sharedConstants.KPIFulfillmentCheckRequestsQueueName, "")
				rabbitMQClient.Dispose()
				if err != nil {
					log.Printf("[MPU][INPUT][consumer=%d] Consumption failed from '%s': %s", consumerID, sharedConstants.KPIFulfillmentCheckRequestsQueueName, err.Error())
				}
				time.Sleep(time.Second)
			}
		}(consumerID)
	}
	wg.Wait()
}

func checkForKPIReprocessRequests() {
	workerCount := getWorkerCount("MPU_REPROCESS_WORKERS", "[MPU][REPROCESS]")
	log.Printf("[MPU][REPROCESS] Starting %d reprocess worker(s)", workerCount)
	var wg sync.WaitGroup
	for workerID := 1; workerID <= workerCount; workerID++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				rabbitMQClient := rabbitmq.NewClient()
				err := rabbitmq.ConsumeJSONMessages[sharedModel.KPIReprocessRequestISCMessage](rabbitMQClient, sharedConstants.KPIReprocessRequestQueueName,
					func(req sharedModel.KPIReprocessRequestISCMessage) error {
						log.Printf("[MPU][REPROCESS][worker=%d] Starting KPI reprocess for KPI definition %s", workerID, req.KPIDefinitionUID)
						return processing.ReprocessKPI(req)
					},
				)
				rabbitMQClient.Dispose()
				if err != nil {
					log.Printf("[MPU][REPROCESS][worker=%d] Failed consuming KPI reprocess queue: %s", workerID, err.Error())
				}
				time.Sleep(time.Second)
			}
		}(workerID)
	}
	wg.Wait()
}

func checkForKPIDefinitionsBySDTypeDenotationMapUpdates() {
	for {
		rabbitMQClient := rabbitmq.NewClient()
		processing.DenotationMapUpdates(rabbitMQClient)
		rabbitMQClient.Dispose()
		time.Sleep(time.Second)
	}
}

func checkForSDTypeUpdates() {
	for {
		rabbitMQClient := rabbitmq.NewClient()
		err := rabbitmq.ConsumeJSONMessages[[]sharedModel.SDTypeUpdateISCMessage](rabbitMQClient, sharedConstants.SetOfSDTypesUpdatesQueueName,
			func(messages []sharedModel.SDTypeUpdateISCMessage) error {
				processing.UpdateSDType(messages)
				dispatch.EnsureIngestConsumersForSDTypes(messages)
				return nil
			},
		)
		rabbitMQClient.Dispose()
		if err != nil {
			log.Printf("[MPU][SDTYPE] Failed consuming SDType updates: %s", err.Error())
		}
		time.Sleep(time.Second)
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
	rawBackendCoreURL := sharedUtils.GetEnvironmentVariableValue("BACKEND_CORE_URL").GetPayloadOrDefault("http://backend-core:9090")
	parsedBackendCoreURL, err := url.Parse(rawBackendCoreURL)
	sharedUtils.TerminateOnError(err, fmt.Sprintf("Unable to parse the backend-core URL: %s", rawBackendCoreURL))
	sharedUtils.TerminateOnError(sharedUtils.WaitForDSs(time.Minute, sharedUtils.NewPairOf(parsedBackendCoreURL.Hostname(), parsedBackendCoreURL.Port())), "Some dependencies of this application are inaccessible")
	log.Println("Dependencies should be up and running...")
	processing.InitializeProcessing()
	dispatch.StartInputProcessingWorkers()
	go dispatch.LogInputStatsPeriodically(time.Minute)
	go dispatch.LogInputTrendsPeriodically(dispatch.InputTrendSampleInterval)
	go checkForKPIDefinitionsBySDTypeDenotationMapUpdates()
	sharedUtils.TerminateOnError(processing.BootstrapProcessingState(), "Unable to bootstrap processing state")
	sharedUtils.StartLoggingProfilingInformationPeriodically(time.Minute)
	sharedUtils.WaitForAll(checkForKPIFulfilmentCheckRequests, checkForKPIReprocessRequests, checkForSDTypeUpdates, checkForKPIDeleteRequests)
	processing.StartCacheCleanup(time.Hour, processing.CacheTTL)
}
