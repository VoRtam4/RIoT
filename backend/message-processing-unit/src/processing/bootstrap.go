package processing

import (
	"errors"
	"log"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var errBootstrapStreamCompleted = errors.New("bootstrap stream completed")

func BootstrapProcessingState() error {
	log.Printf("[MPU][BOOTSTRAP] Starting runtime cache bootstrap | unit=%s", unitUUID)

	rawReady := make(chan error, 1)
	kpiReady := make(chan error, 1)

	go func() {
		rawReady <- bootstrapRawCache()
	}()
	go func() {
		kpiReady <- bootstrapKPICache()
	}()

	time.Sleep(time.Second)

	if err := notifyBackendCoreOfConnection(); err != nil {
		return err
	}

	<-initialConfigReady

	if err := <-rawReady; err != nil {
		return err
	}
	if err := <-kpiReady; err != nil {
		return err
	}

	log.Printf("[MPU][BOOTSTRAP] Runtime cache bootstrap finished | unit=%s", unitUUID)
	return nil
}

func notifyBackendCoreOfConnection() error {
	client := rabbitmq.NewClient()
	defer client.Dispose()

	return client.PublishJSONMessage(
		sharedUtils.NewEmptyOptional[string](),
		sharedUtils.NewOptionalOf(sharedConstants.MessageProcessingUnitConnectionNotificationsQueueName),
		[]byte("{}"),
	)
}

func bootstrapRawCache() error {
	client := rabbitmq.NewClient()
	defer client.Dispose()

	entries := 0
	skipped := 0
	err := rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.RawDataPointCacheBootstrapISCMessage](
		client,
		sharedConstants.RawDataPointCacheBootstrapQueueName,
		"",
		func(messagePayload sharedModel.RawDataPointCacheBootstrapISCMessage, delivery amqp.Delivery) error {
			now := time.Now()
			for _, msg := range messagePayload.Tuple {
				if now.Sub(msg.EventTime) > CacheTTL {
					skipped++
					continue
				}
				deserializationResult := sharedUtils.DeserializeFromJSON[map[string]interface{}](msg.Payload)
				if deserializationResult.IsFailure() {
					return deserializationResult.GetError()
				}
				lastRaw.Store(msg.SDInstanceUID, sharedModel.RawState{
					Values:         deserializationResult.GetPayload(),
					EventTime:      msg.EventTime,
					SynchronizedAt: now,
				})
				entries++
			}
			if messagePayload.Done {
				if err := delivery.Ack(false); err != nil {
					return err
				}
				log.Printf("[MPU][BOOTSTRAP][RAW] Received bootstrap completion | entries=%d skipped=%d", entries, skipped)
				return errBootstrapStreamCompleted
			}
			return nil
		},
	)
	if err != nil && !errors.Is(err, errBootstrapStreamCompleted) {
		return err
	}
	return nil
}

func bootstrapKPICache() error {
	client := rabbitmq.NewClient()
	defer client.Dispose()

	entries := 0
	skipped := 0
	err := rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.KPIFulfillmentCacheBootstrapISCMessage](
		client,
		sharedConstants.KPIFulfillmentCacheBootstrapQueueName,
		"",
		func(messagePayload sharedModel.KPIFulfillmentCacheBootstrapISCMessage, delivery amqp.Delivery) error {
			now := time.Now()
			for _, msg := range messagePayload.Tuple {
				if now.Sub(msg.EventTime) > CacheTTL {
					skipped++
					continue
				}
				lastKPI.Store(sharedModel.KPIKey{
					SDInstanceUID:   msg.SDInstanceUID,
					KPIDefinitionID: msg.KPIDefinitionID,
				}, sharedModel.KPIState{
					Value:          msg.Fulfilled,
					EventTime:      msg.EventTime,
					SynchronizedAt: now,
				})
				entries++
			}
			if messagePayload.Done {
				if err := delivery.Ack(false); err != nil {
					return err
				}
				log.Printf("[MPU][BOOTSTRAP][KPI] Received bootstrap completion | entries=%d skipped=%d", entries, skipped)
				return errBootstrapStreamCompleted
			}
			return nil
		},
	)
	if err != nil && !errors.Is(err, errBootstrapStreamCompleted) {
		return err
	}
	return nil
}
