/**
 * @file consumer.go
 * @brief RabbitMQ delivery loop, interní input buffer a processing workery dispatch vrstvy MPU.
 *
 * @author Vojtěch Hubáček
 *
 * @ingroup riot_message_processing_unit
 */
package dispatch

import (
	"fmt"
	"log"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/MichalBures-OG/bp-bures-RIoT-message-processing-unit/src/processing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func StartInputProcessingWorkers() {
	inputWorkersOnce.Do(func() {
		workerCount := getPositiveInt("MPU_INPUT_WORKERS", 1, "[MPU][INPUT]")
		queueCapacity := getPositiveInt("MPU_INTERNAL_INPUT_BUFFER", 1000, "[MPU][INPUT]")
		configureInputScheduler("[MPU][INPUT]")
		processing.SetInputConcurrency(workerCount)
		inputProcessingQueue = make(chan inputProcessingTask, queueCapacity)
		log.Printf("[MPU][INPUT] Starting %d input worker(s) | internalBuffer=%d", workerCount, queueCapacity)
		for workerID := 1; workerID <= workerCount; workerID++ {
			go runInputProcessingWorker(workerID)
		}
	})
}

func runInputProcessingWorker(workerID int) {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	for task := range inputProcessingQueue {
		err := processing.ProcessKPIFulfillmentCheckRequestTuple(rabbitMQClient, task.tuple)
		stats := getInputSourceStats(task.source)
		stats.enqueued.Add(-1)
		if err != nil {
			stats.failed.Add(1)
			log.Printf("[MPU][INPUT][worker=%d] Failed processing input from %s: %s", workerID, task.source, err.Error())
		} else {
			stats.processed.Add(1)
		}
		task.result <- err
	}
}

func processInputTupleViaInternalQueue(source string, tuple sharedModel.KPIFulfillmentCheckRequestTupleISCMessage) error {
	if len(tuple) == 0 {
		return nil
	}
	waitForInputCapacity(source)
	stats := getInputSourceStats(source)
	result := make(chan error, 1)
	stats.received.Add(1)
	stats.enqueued.Add(1)
	inputProcessingQueue <- inputProcessingTask{
		source: source,
		tuple:  tuple,
		result: result,
	}
	return <-result
}

func ConsumeInputQueue(rabbitMQClient rabbitmq.Client, queueName string, consumerName string) error {
	deliveries, err := rabbitMQClient.Consume(queueName, consumerName)
	if err != nil {
		return err
	}
	for delivery := range deliveries {
		if err := processInputDelivery(queueName, delivery); err != nil {
			return err
		}
	}
	return fmt.Errorf("message channel closed for queue %s", queueName)
}

func processInputDelivery(source string, delivery amqp.Delivery) error {
	stats := getInputSourceStats(source)
	if delivery.ContentType != "application/json" {
		stats.failed.Add(1)
		if err := delivery.Nack(false, false); err != nil {
			return fmt.Errorf("failed to reject non-json message from %s: %w", source, err)
		}
		return fmt.Errorf("incorrect message content type from %s: %s", source, delivery.ContentType)
	}
	jsonDeserializationResult := sharedUtils.DeserializeFromJSON[sharedModel.KPIFulfillmentCheckRequestTupleISCMessage](delivery.Body)
	if jsonDeserializationResult.IsFailure() {
		stats.failed.Add(1)
		if err := delivery.Nack(false, false); err != nil {
			return fmt.Errorf("failed to reject invalid json message from %s: %w", source, err)
		}
		return fmt.Errorf("failed to deserialize input message from %s: %w", source, jsonDeserializationResult.GetError())
	}
	if isInputSourceDropReady(stats) {
		if inputDropperEnabled {
			stats.dropped.Add(1)
			if err := delivery.Ack(false); err != nil {
				return fmt.Errorf("failed to ack dropped message from %s: %w", source, err)
			}
			log.Printf("[MPU][INPUT][DROPPER] Dropped message | source=%s dropped=%d", source, stats.dropped.Load())
			return nil
		}
		stats.wouldDrop.Add(1)
		log.Printf("[MPU][INPUT][DROPPER] Dry-run would drop message | source=%s wouldDrop=%d", source, stats.wouldDrop.Load())
	}
	if err := processInputTupleViaInternalQueue(source, jsonDeserializationResult.GetPayload()); err != nil {
		if nackErr := delivery.Nack(false, true); nackErr != nil {
			return fmt.Errorf("failed to requeue message from %s after processing error %w: %w", source, err, nackErr)
		}
		return err
	}
	if err := delivery.Ack(false); err != nil {
		return fmt.Errorf("failed to ack message from %s: %w", source, err)
	}
	return nil
}
