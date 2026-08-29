/**
 * @file ingest.go
 * @brief Per-SDType ingest fronty a konzumenti dispatch vrstvy MPU.
 *
 * @author Vojtěch Hubáček
 *
 * @ingroup riot_message_processing_unit
 */
package dispatch

import (
	"fmt"
	"log"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

func ensureIngestQueueForSDType(rabbitMQClient rabbitmq.Client, sdTypeUID string) error {
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

func EnsureIngestConsumersForSDTypes(messages []sharedModel.SDTypeUpdateISCMessage) {
	for _, message := range messages {
		ensureIngestConsumerForSDType(message.SDTypeUID)
	}
}

func ensureIngestConsumerForSDType(sdTypeUID string) {
	if sdTypeUID == "" {
		return
	}
	if _, loaded := ingestConsumerRegistry.LoadOrStore(sdTypeUID, struct{}{}); loaded {
		return
	}
	go runIngestConsumer(sdTypeUID)
}

func runIngestConsumer(sdTypeUID string) {
	queueName := sharedConstants.IngestQueueName(sdTypeUID)
	log.Printf("[MPU][INGEST] Starting consumer for SDType %s queue=%s", sdTypeUID, queueName)
	for {
		rabbitMQClient := rabbitmq.NewClient()
		if err := ensureIngestQueueForSDType(rabbitMQClient, sdTypeUID); err != nil {
			rabbitMQClient.Dispose()
			log.Printf("[MPU][INGEST] Failed ensuring queue for SDType %s: %s", sdTypeUID, err.Error())
			time.Sleep(time.Second)
			continue
		}
		err := ConsumeInputQueue(rabbitMQClient, queueName, "")
		rabbitMQClient.Dispose()
		if err != nil {
			log.Printf("[MPU][INGEST] Consumption failed from '%s': %s", queueName, err.Error())
		}
		time.Sleep(time.Second)
	}
}
