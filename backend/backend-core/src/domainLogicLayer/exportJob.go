package domainLogicLayer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type ExportJob struct {
	ID       string
	FilePath string
	Done     chan struct{}
	Err      error
}

var (
	exportJobs   = map[string]*ExportJob{}
	exportJobsMu sync.RWMutex
)

func StartTimeSeriesExportAggregateKPI(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput) (string, error) {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, graphQLModel.TimeSeriesTypeKpi, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return "", err
	}
	req := gql2dll.ToDLLTimeSeriesReadKPIRequest(input, sdTypeUID, instanceUIDs)
	return StartTimeSeriesExportBase(userID, req, input.SdTypeID, sdTypeUID, instanceUIDs)
}

func StartTimeSeriesExport(userID uint32, input graphQLModel.TimeSeriesReadInput) (string, error) {
	sdTypeUID, instanceUIDs, err := resolveAndValidate(userID, input.Type, input.SdTypeID, input.SdInstanceIDs, input.KpiDefinitionIDs)
	if err != nil {
		return "", err
	}
	req := gql2dll.ToDLLTimeSeriesReadRequest(input, sdTypeUID, instanceUIDs)
	return StartTimeSeriesExportBase(userID, req, input.SdTypeID, sdTypeUID, instanceUIDs)
}

func StartTimeSeriesExportBase(userID uint32, req sharedModel.TimeSeriesReadRequest, sdTypeID *uint32, sdTypeUID string, instanceUIDs []string) (string, error) {
	id := uuid.New().String()
	job := &ExportJob{
		ID:       id,
		FilePath: "/tmp/export_" + id + ".csv",
		Done:     make(chan struct{}),
	}
	exportJobsMu.Lock()
	exportJobs[id] = job
	exportJobsMu.Unlock()
	go runExportJob(job, req, sdTypeID, sdTypeUID, instanceUIDs)
	return id, nil
}

func runExportJob(job *ExportJob, req sharedModel.TimeSeriesReadRequest, sdTypeID *uint32, sdTypeUID string, instanceUIDs []string) {
	defer close(job.Done)
	req.Limit = nil
	req.Batch = nil
	client := rabbitmq.NewClient()
	defer client.Dispose()
	ch := client.GetChannel()
	replyQueue, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		job.Err = err
		return
	}
	msgs, err := ch.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		job.Err = err
		return
	}
	handler := buildExportHandler(job, req, sdTypeID)
	if handler == nil {
		return
	}
	correlationID := job.ID
	jsonData, err := json.Marshal(req)
	if err != nil {
		job.Err = err
		return
	}
	err = ch.PublishWithContext(context.Background(), "", sharedConstants.TimeSeriesReadRequestQueueName, false, false, amqp091.Publishing{
		ContentType:   "application/json",
		Body:          jsonData,
		CorrelationId: correlationID,
		ReplyTo:       replyQueue.Name,
	})
	if err != nil {
		job.Err = err
		return
	}
	err = rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesReadResponse](msgs, correlationID, 0, func(resp sharedModel.TimeSeriesReadResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			job.Err = fmt.Errorf("%s", resp.Error)
			return true, job.Err
		}
		if err := handler(resp, msg); err != nil {
			job.Err = err
			return true, err
		}
		if !resp.HasMoreBatches {
			go cleanupJob(job.ID)
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		job.Err = err
	}
}

func buildExportHandler(job *ExportJob, req sharedModel.TimeSeriesReadRequest, sdTypeID *uint32) func(resp sharedModel.TimeSeriesReadResponse, _ amqp091.Delivery) error {
	file, err := os.Create(job.FilePath)
	if err != nil {
		job.Err = err
		return nil
	}
	writer := csv.NewWriter(file)
	params := loadParametersFromDB(graphQLModel.TimeSeriesType(req.Type), sdTypeID)
	headerWritten := false
	return func(resp sharedModel.TimeSeriesReadResponse, _ amqp091.Delivery) error {
		defer func() {
			if !resp.HasMoreBatches {
				writer.Flush()
				file.Close()
			}
		}()
		if resp.Error != "" {
			job.Err = fmt.Errorf("%s", resp.Error)
			return job.Err
		}
		if !headerWritten {
			header := make([]string, 0, len(params))
			for _, p := range params {
				header = append(header, p.Denotation)
			}
			_ = writer.Write(header)
			headerWritten = true
		}
		for _, point := range resp.Data {
			row := make([]string, 0, len(params))
			for _, param := range params {
				switch param.Role {

				case graphQLModel.ParameterRoleTime:
					row = append(row, point.Time.Format(time.RFC3339Nano))

				case graphQLModel.ParameterRoleMeta:
					val := ""
					if resp.Base != nil {
						val = resp.Base[param.Denotation]
					}
					if val == "" {
						if v, ok := point.Tags[param.Denotation]; ok {
							val = v
						}
					}
					if val == "" && param.Denotation == "kpiDefinitionID" && len(req.KPIDefinitionIDs) > 0 {
						val = fmt.Sprint(req.KPIDefinitionIDs[0])
					}
					row = append(row, val)

				case graphQLModel.ParameterRoleTag:
					if val, ok := point.Tags[param.Denotation]; ok {
						row = append(row, val)
					} else {
						row = append(row, "")
					}

				case graphQLModel.ParameterRoleField:
					if val, ok := point.Data[param.Denotation]; ok {
						row = append(row, fmt.Sprint(val))
					} else {
						row = append(row, "")
					}
				default:
					row = append(row, "")
				}
			}
			_ = writer.Write(row)
		}
		writer.Flush()
		return nil
	}
}

func GetExportJob(id string) (*ExportJob, bool) {
	exportJobsMu.RLock()
	defer exportJobsMu.RUnlock()
	job, ok := exportJobs[id]
	return job, ok
}

func cleanupJob(id string) {
	time.Sleep(30 * time.Minute)
	exportJobsMu.Lock()
	job := exportJobs[id]
	delete(exportJobs, id)
	exportJobsMu.Unlock()
	if job != nil {
		_ = os.Remove(job.FilePath)
	}
}
