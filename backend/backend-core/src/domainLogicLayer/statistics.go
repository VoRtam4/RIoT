package domainLogicLayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	amqp "github.com/rabbitmq/amqp091-go"
)

func randomString(l int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func Query(input sharedModel.TimeSeriesReadRequest) sharedUtils.Result[[]graphQLModel.OutputData] {

	request := sharedUtils.SerializeToJSON(input)

	rabbitMQClient := getDLLRabbitMQClient()

	correlationId := randomString(32)

	outputChannel := make(chan sharedUtils.Result[[]sharedModel.TimeSeriesDataPoint])

	go func() {

		client := rabbitmq.NewClient()
		defer client.Dispose()

		err := rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesReadResponse](
			client,
			sharedConstants.TimeSeriesReadResponseQueueName,
			correlationId,
			func(resp sharedModel.TimeSeriesReadResponse, delivery amqp.Delivery) error {

				if resp.Error != "" {
					outputChannel <- sharedUtils.NewFailureResult[[]sharedModel.TimeSeriesDataPoint](errors.New(resp.Error))
				} else {
					outputChannel <- sharedUtils.NewSuccessResult(resp.Data)
				}

				close(outputChannel)
				return nil
			},
		)

		if err != nil {
			log.Printf("Statistics Query | %s", err)
			outputChannel <- sharedUtils.NewFailureResult[[]sharedModel.TimeSeriesDataPoint](err)
			close(outputChannel)
		}
	}()

	err := rabbitMQClient.PublishJSONMessageRPC(
		sharedUtils.NewEmptyOptional[string](),
		sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesReadRequestQueueName),
		request.GetPayload(),
		correlationId,
		sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesReadResponseQueueName),
	)

	if err != nil {
		log.Printf("Statistics Query | %s", err)
		return sharedUtils.NewFailureResult[[]graphQLModel.OutputData](err)
	}

	result := <-outputChannel

	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.OutputData](result.GetError())
	}

	convertedResult, err := ConvertOutputData(result.GetPayload())

	if err != nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.OutputData](err)
	}

	return sharedUtils.NewSuccessResult(convertedResult)
}

func Save(input graphQLModel.InputData) sharedUtils.Result[bool] {

	// Save je nyní zbytečné, protože RAW data zapisuje MPU.
	// Funkci ale necháme kvůli kompatibilitě.

	log.Println("Statistics Save ignored (handled by MPU)")
	return sharedUtils.NewSuccessResult(true)
}

func MapStatisticsInputToReadRequestBody(
	statsInput *graphQLModel.StatisticsInput,
	simpleSensors *graphQLModel.SimpleSensors,
	sensorsWithFields *graphQLModel.SensorsWithFields,
) (*sharedModel.TimeSeriesReadRequest, error) {

	req := &sharedModel.TimeSeriesReadRequest{}

	if simpleSensors != nil {
		for k := range simpleSensors.Sensors {
			s := strconv.Itoa(k)
			req.SDInstanceUID = &s
			break
		}
	}

	if sensorsWithFields != nil {
		for _, sensor := range sensorsWithFields.Sensors {
			req.SDInstanceUID = &sensor.Key
			break
		}
	}

	if statsInput != nil {

		if statsInput.AggregateMinutes != nil {
			req.AggregateMinutes = statsInput.AggregateMinutes
		}

		if statsInput.From != nil {
			t, err := time.Parse(time.RFC3339, *statsInput.From)
			if err != nil {
				return nil, err
			}
			req.From = &t
		}

		if statsInput.To != nil {
			t, err := time.Parse(time.RFC3339, *statsInput.To)
			if err != nil {
				return nil, err
			}
			req.To = &t
		}
	}

	return req, nil
}

func ConvertOutputData(sharedData []sharedModel.TimeSeriesDataPoint) ([]graphQLModel.OutputData, error) {

	var result []graphQLModel.OutputData

	for _, item := range sharedData {

		timeString := item.Time.Format(time.RFC3339)

		dataBytes, err := json.Marshal(item.Data)
		if err != nil {
			return nil, fmt.Errorf("error marshaling data: %v", err)
		}

		dataString := string(dataBytes)

		deviceID := ""
		if v, ok := item.Tags["sdInstanceUID"]; ok {
			deviceID = v
		}

		var deviceType *string
		if v, ok := item.Tags["sdType"]; ok {
			deviceType = &v
		}

		result = append(result, graphQLModel.OutputData{
			Time:       timeString,
			DeviceID:   deviceID,
			DeviceType: deviceType,
			Data:       dataString,
		})
	}

	return result, nil
}
