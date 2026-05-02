package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
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
	defer influx.Close()
	log.Println("Time Series Store Ready")
	sharedUtils.WaitForAll(func() {
		runConsumerLoop("raw records", func() error { return consumeRawRecords(influx) })
	}, func() {
		runConsumerLoop("kpi records", func() error { return consumeKPIRecords(influx) })
	}, func() {
		runConsumerLoop("read requests", func() error { return consumeReadRequests(influx) })
	}, func() {
		runConsumerLoop("distinct tag value requests", func() error { return consumeDistinctTagValueRequests(influx) })
	}, func() {
		runConsumerLoop("reprocess read requests", func() error { return consumeReprocessReadRequests(influx) })
	}, func() {
		runConsumerLoop("delete requests", func() error { return consumeDeleteRequests(influx) })
	})
}

func runConsumerLoop(name string, consumer func() error) {
	for {
		if err := consumer(); err != nil {
			log.Printf("[TSS] Consumer '%s' stopped: %s", name, err.Error())
		}
		time.Sleep(time.Second)
	}
}

func consumeRawRecords(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessages[[]sharedModel.TimeSeriesRawRecord](
		rabbitMQClient, sharedConstants.TimeSeriesRawDataQueueName,
		func(records []sharedModel.TimeSeriesRawRecord) error {
			influx.WriteRawBatch(records)
			return nil
		},
	)
}

func consumeKPIRecords(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessages[[]sharedModel.TimeSeriesKPIResultRecord](
		rabbitMQClient, sharedConstants.TimeSeriesKPIResultQueueName,
		func(records []sharedModel.TimeSeriesKPIResultRecord) error {
			influx.WriteKPIBatch(records)
			return nil
		},
	)
}

func consumeReadRequests(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesReadRequest](rabbitMQClient, sharedConstants.TimeSeriesReadRequestQueueName, "",
		func(req sharedModel.TimeSeriesReadRequest, delivery amqp.Delivery) error {
			go handleReadRequest(influx, req, delivery)
			return nil
		},
	)
}

func handleReadRequest(influx internal.Influx2Client, req sharedModel.TimeSeriesReadRequest, delivery amqp.Delivery) {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()

	base := map[string]string{
		"type":      string(req.Type),
		"sdTypeUID": req.SDTypeUID,
	}
	plan := internal.BuildQueryPlan(req)
	var lastPoint *sharedModel.TimeSeriesDataPoint
	buildCursor := func(p *sharedModel.TimeSeriesDataPoint) *sharedModel.TimeSeriesCursor {
		if p == nil {
			return nil
		}
		cursor := &sharedModel.TimeSeriesCursor{
			Time:          p.Time,
			SDInstanceUID: p.Tags["sdInstanceUID"],
		}
		if plan.IsKPI {
			if v, ok := p.Tags["kpiDefinitionID"]; ok {
				if parsed, err := strconv.ParseUint(v, 10, 32); err == nil {
					kpi := uint32(parsed)
					cursor.KPIDefinitionID = &kpi
				}
			}
		}
		return cursor
	}
	ch := rabbitMQClient.GetChannel()
	publish := func(resp sharedModel.TimeSeriesReadResponse) error {
		jsonData, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		return ch.PublishWithContext(context.Background(), "", delivery.ReplyTo, false, false, amqp.Publishing{
			ContentType:   "application/json",
			Body:          jsonData,
			CorrelationId: delivery.CorrelationId,
		})
	}
	send := func(rows []sharedModel.TimeSeriesDataPoint, hasMoreBatches bool, hasMoreData bool) error {
		if len(rows) > 0 {
			p := rows[len(rows)-1]
			lastPoint = &p
		}
		resp := sharedModel.TimeSeriesReadResponse{
			Base:           base,
			Data:           rows,
			HasMoreBatches: hasMoreBatches,
			HasMoreData:    hasMoreData,
		}
		if !hasMoreBatches {
			cursor := buildCursor(lastPoint)
			if cursor == nil && req.Cursor != nil {
				cursor = req.Cursor
			}
			resp.NextCursor = cursor
		}
		return publish(resp)
	}
	err := influx.StreamRead(plan, send)
	if err != nil {
		resp := sharedModel.TimeSeriesReadResponse{
			Base:  base,
			Error: err.Error(),
		}
		if publishErr := publish(resp); publishErr != nil {
			log.Printf("[TSS] failed to publish read error response: %s", publishErr.Error())
		}
	}
}

func consumeDistinctTagValueRequests(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesDistinctTagValuesRequest](
		rabbitMQClient,
		sharedConstants.TimeSeriesDistinctTagValuesRequestQueueName,
		"",
		func(req sharedModel.TimeSeriesDistinctTagValuesRequest, delivery amqp.Delivery) error {
			ch := rabbitMQClient.GetChannel()
			publish := func(resp sharedModel.TimeSeriesDistinctTagValuesResponse) error {
				jsonData, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				return ch.PublishWithContext(context.Background(), "", delivery.ReplyTo, false, false, amqp.Publishing{
					ContentType:   "application/json",
					Body:          jsonData,
					CorrelationId: delivery.CorrelationId,
				})
			}

			values, err := influx.DistinctTagValues(req)
			if err != nil {
				return publish(sharedModel.TimeSeriesDistinctTagValuesResponse{
					Error: err.Error(),
				})
			}

			return publish(sharedModel.TimeSeriesDistinctTagValuesResponse{
				Values: values,
			})
		},
	)
}

func consumeReprocessReadRequests(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesReprocessReadRequest](rabbitMQClient, sharedConstants.TimeSeriesReprocessReadRequestQueueName, "",
		func(req sharedModel.TimeSeriesReprocessReadRequest, delivery amqp.Delivery) error {
			ch := rabbitMQClient.GetChannel()
			send := func(resp sharedModel.TimeSeriesReprocessReadResponse) error {
				if delivery.ReplyTo == "" {
					return fmt.Errorf("missing ReplyTo in request")
				}
				jsonData, err := json.Marshal(resp)
				if err != nil {
					return err
				}
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				return ch.PublishWithContext(ctx, "", delivery.ReplyTo, false, false, amqp.Publishing{
					ContentType:   "application/json",
					Body:          jsonData,
					CorrelationId: delivery.CorrelationId,
				})
			}
			err := influx.StreamReprocess(req, func(points []sharedModel.TimeSeriesDataPoint, hasMore bool) error {
				return send(sharedModel.TimeSeriesReprocessReadResponse{
					Data:    points,
					HasMore: hasMore,
				})
			})
			if err != nil {
				return send(sharedModel.TimeSeriesReprocessReadResponse{
					Error: err.Error(),
				})
			}
			return send(sharedModel.TimeSeriesReprocessReadResponse{
				HasMore: false,
			})
		},
	)
}

func consumeDeleteRequests(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessages[sharedModel.KPIDeleteResultsRequestISCMessage](rabbitMQClient, sharedConstants.TSDBDeleteQueueName, influx.DeleteKPI)
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
