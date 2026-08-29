/**
 * @file main.go
 * @brief Vstupní bod modulu Time Series Store platformy RIoT a registrace jeho komunikačních toků.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @defgroup riot_time_series_store Time Series Store
 * @ingroup riot
 * @see ../README.md
 */
package main

import (
	"context"
	"encoding/json"
	"errors"
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
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xjohnp00/jiap/backend/shared/time-series-store/src/internal"
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

func main() {
	hasError, environment := parseParameters()
	if hasError {
		os.Exit(1)
	}
	log.Println("Waiting for dependencies...")
	rawBackendCoreURL := sharedUtils.GetEnvironmentVariableValue("BACKEND_CORE_URL").GetPayloadOrDefault("http://backend-core:9090")
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
		runConsumerLoop("read cancel requests", func() error { return consumeReadCancelRequests() })
	}, func() {
		runConsumerLoop("distinct tag value requests", func() error { return consumeDistinctTagValueRequests(influx) })
	}, func() {
		runConsumerLoop("reprocess read requests", func() error { return consumeReprocessReadRequests(influx) })
	}, func() {
		runConsumerLoop("record neighborhood requests", func() error { return consumeRecordNeighborhoodRequests(influx) })
	}, func() {
		runConsumerLoop("late record corrections", func() error { return consumeLateRecordCorrections(influx) })
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
	workerCount := getWorkerCount("INFLUX_READ_WORKERS", "[TSS][READ]")
	log.Printf("[TSS][READ] Starting %d read worker(s)", workerCount)
	var wg sync.WaitGroup
	for workerID := 1; workerID <= workerCount; workerID++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				if err := consumeReadRequestsWorker(influx); err != nil {
					log.Printf("[TSS][READ][worker=%d] Consumer stopped: %s", workerID, err.Error())
				}
				time.Sleep(time.Second)
			}
		}(workerID)
	}
	wg.Wait()
	return nil
}

func consumeReadRequestsWorker(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesReadRequest](rabbitMQClient, sharedConstants.TimeSeriesReadRequestQueueName, "",
		func(req sharedModel.TimeSeriesReadRequest, delivery amqp.Delivery) error {
			handleReadRequest(influx, req, delivery)
			return nil
		},
	)
}

func handleReadRequest(influx internal.Influx2Client, req sharedModel.TimeSeriesReadRequest, delivery amqp.Delivery) {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	internal.RegisterReadJob(req.JobID, cancel)
	defer internal.UnregisterReadJob(req.JobID)

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
			if v, ok := p.Tags["kpiDefinitionUID"]; ok {
				cursor.KPIDefinitionUID = &v
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
	err := influx.StreamRead(ctx, plan, send)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			resp := sharedModel.TimeSeriesReadResponse{
				Base:           base,
				HasMoreBatches: false,
				HasMoreData:    false,
			}
			if publishErr := publish(resp); publishErr != nil {
				log.Printf("[TSS] failed to publish cancelled read response: %s", publishErr.Error())
			}
			return
		}
		resp := sharedModel.TimeSeriesReadResponse{
			Base:  base,
			Error: err.Error(),
		}
		if publishErr := publish(resp); publishErr != nil {
			log.Printf("[TSS] failed to publish read error response: %s", publishErr.Error())
		}
	}
}

func consumeReadCancelRequests() error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessages[sharedModel.TimeSeriesReadCancelRequest](
		rabbitMQClient,
		sharedConstants.TimeSeriesReadCancelRequestQueueName,
		func(req sharedModel.TimeSeriesReadCancelRequest) error {
			internal.CancelReadJob(req.JobID)
			return nil
		},
	)
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
	workerCount := getWorkerCount("INFLUX_REPROCESS_WORKERS", "[TSS][REPROCESS]")
	log.Printf("[TSS][REPROCESS] Starting %d reprocess read worker(s)", workerCount)
	var wg sync.WaitGroup
	for workerID := 1; workerID <= workerCount; workerID++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				if err := consumeReprocessReadRequestsWorker(influx); err != nil {
					log.Printf("[TSS][REPROCESS][worker=%d] Consumer stopped: %s", workerID, err.Error())
				}
				time.Sleep(time.Second)
			}
		}(workerID)
	}
	wg.Wait()
	return nil
}

func consumeReprocessReadRequestsWorker(influx internal.Influx2Client) error {
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

func consumeRecordNeighborhoodRequests(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesRecordNeighborhoodRequest](
		rabbitMQClient,
		sharedConstants.TimeSeriesRecordNeighborhoodRequestQueueName,
		"",
		func(req sharedModel.TimeSeriesRecordNeighborhoodRequest, delivery amqp.Delivery) error {
			ch := rabbitMQClient.GetChannel()
			send := func(resp sharedModel.TimeSeriesRecordNeighborhoodResponse) error {
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
			response, err := influx.GetRecordNeighborhood(req)
			if err != nil {
				return send(sharedModel.TimeSeriesRecordNeighborhoodResponse{Error: err.Error()})
			}
			return send(response)
		},
	)
}

func consumeLateRecordCorrections(influx internal.Influx2Client) error {
	rabbitMQClient := rabbitmq.NewClient()
	defer rabbitMQClient.Dispose()
	return rabbitmq.ConsumeJSONMessages[sharedModel.TimeSeriesLateRecordCorrection](
		rabbitMQClient,
		sharedConstants.TimeSeriesLateRecordCorrectionQueueName,
		influx.ApplyLateRecordCorrection,
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
