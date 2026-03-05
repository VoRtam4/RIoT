package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xjohnp00/jiap/backend/shared/time-series-store/src/internal"
)

func main() {
	hasError, environment := parseParameters()

	if hasError {
		os.Exit(1)
	}

	log.Println("Waiting for dependencies...")
	rawBackendCoreURL := sharedUtils.GetEnvironmentVariableValue("BACKEND_CORE_URL").GetPayloadOrDefault("http://riot-backend-core:9090")
	parsedBackendCoreURL, err := url.Parse(rawBackendCoreURL)
	sharedUtils.TerminateOnError(err, fmt.Sprintf("Unable to parse the backend-core URL: %s", rawBackendCoreURL))

	parsedInfluxURL, err := url.Parse(environment.InfluxUrl)
	sharedUtils.TerminateOnError(err, fmt.Sprintf("Unable to parse the InfluxDB URL: %s", environment.InfluxUrl))

	parsedRabbitMQURL, err := url.Parse(environment.AmqpURLValue)
	sharedUtils.TerminateOnError(err, fmt.Sprintf("Unable to parse the RabbitMQ URL: %s", environment.AmqpURLValue))

	sharedUtils.TerminateOnError(sharedUtils.WaitForDSs(time.Minute,
		sharedUtils.NewPairOf(parsedBackendCoreURL.Hostname(), parsedBackendCoreURL.Port()),
		sharedUtils.NewPairOf(parsedInfluxURL.Hostname(), parsedInfluxURL.Port()),
		sharedUtils.NewPairOf(parsedRabbitMQURL.Hostname(), parsedRabbitMQURL.Port()),
	), "Some dependencies of this application are inaccessible")
	log.Println("Dependencies should be up and running...")

	log.Println("Time Series Store Starting")

	influx := internal.NewInflux2Client(environment.InfluxUrl, environment.InfluxToken, environment.InfluxOrg, environment.InfluxBucket)

	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	defer influx.Close()

	log.Println("Time Series Store Ready")

	sharedUtils.WaitForAll(
		func() {
			if err := consumeRawRecords(rabbitMQClient, influx); err != nil {
				log.Println(err.Error())
			}
		},
		func() {
			if err := consumeKPIRecords(rabbitMQClient, influx); err != nil {
				log.Println(err.Error())
			}
		},
		func() {
			if err := consumeReadRequests(rabbitMQClient, influx); err != nil {
				log.Println(err.Error())
			}
		},
	)
}

func consumeRawRecords(rabbitMQClient rabbitmq.Client, influx internal.Influx2Client) error {
	return rabbitmq.ConsumeJSONMessages[sharedModel.TimeSeriesRawRecord](
		rabbitMQClient,
		sharedConstants.TimeSeriesRawDataQueueName,
		func(record sharedModel.TimeSeriesRawRecord) error {
			influx.WriteRaw(record)
			return nil
		},
	)
}

func consumeKPIRecords(rabbitMQClient rabbitmq.Client, influx internal.Influx2Client) error {
	return rabbitmq.ConsumeJSONMessages[sharedModel.TimeSeriesKPIResultRecord](
		rabbitMQClient,
		sharedConstants.TimeSeriesKPIResultQueueName,
		func(record sharedModel.TimeSeriesKPIResultRecord) error {
			influx.WriteKPI(record)
			return nil
		},
	)
}

func consumeReadRequests(rabbitMQClient rabbitmq.Client, influx internal.Influx2Client) error {
	return rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesReadRequest](
		rabbitMQClient,
		sharedConstants.TimeSeriesReadRequestQueueName,
		"",
		func(req sharedModel.TimeSeriesReadRequest, delivery amqp.Delivery) error {
			result := influx.Query(req)

			resp := sharedModel.TimeSeriesReadResponse{
				Data:  nil,
				Error: "",
			}

			if result.IsSuccess() {
				resp.Data = result.GetPayload()
			} else {
				resp.Error = result.GetError().Error()
			}

			jsonData, err := json.Marshal(resp)
			if err != nil {
				log.Printf("Error during marshal: %s", err)
				return nil
			}

			err = rabbitMQClient.PublishJSONMessageRPC(
				sharedUtils.NewEmptyOptional[string](),
				sharedUtils.NewOptionalOf(delivery.ReplyTo),
				jsonData,
				delivery.CorrelationId,
				sharedUtils.NewEmptyOptional[string](),
			)

			if err != nil {
				log.Printf("Error publishing RPC response: %s", err)
				return err
			}
			return nil
		},
	)
}

func parseParameters() (bool, internal.TimeSeriesStoreEnvironment) {
	token := sharedUtils.GetEnvironmentVariableValue("INFLUX_TOKEN")
	url := sharedUtils.GetEnvironmentVariableValue("INFLUX_URL")
	org := sharedUtils.GetEnvironmentVariableValue("INFLUX_ORGANIZATION")
	bucket := sharedUtils.GetEnvironmentVariableValue("INFLUX_BUCKET")
	ampqUrl := sharedUtils.GetEnvironmentVariableValue("RABBITMQ_URL")

	hasError := token.IsEmpty() || url.IsEmpty() || org.IsEmpty() || bucket.IsEmpty() || ampqUrl.IsEmpty()

	if token.IsEmpty() {
		log.Fatalln("Empty token")
	}

	if url.IsEmpty() {
		log.Fatalln("Empty url")
	}

	if org.IsEmpty() {
		log.Fatalln("Empty org")
	}

	if bucket.IsEmpty() {
		log.Fatalln("Empty bucket")
	}

	if ampqUrl.IsEmpty() {
		log.Fatalln("Empty ampqUrl")
	}

	environment := sharedUtils.Ternary[internal.TimeSeriesStoreEnvironment](!hasError, internal.TimeSeriesStoreEnvironment{
		InfluxToken:  token.GetPayload(),
		InfluxUrl:    url.GetPayload(),
		InfluxOrg:    org.GetPayload(),
		InfluxBucket: bucket.GetPayload(),
		AmqpURLValue: ampqUrl.GetPayload(),
	}, internal.TimeSeriesStoreEnvironment{})

	return hasError, environment
}
