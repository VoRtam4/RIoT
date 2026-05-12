package domainLogicLayer

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/events"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/rabbitmq/amqp091-go"
)

type ExportJob struct {
	ID               uint32
	UserID           uint32
	FilePath         string
	Done             chan struct{}
	Err              error
	Status           graphQLModel.ExportStatus
	CreatedAt        time.Time
	ExpiresAt        *time.Time
	DownloadURL      *string
	ErrorMessage     *string
	CancelRequested  bool
	RequestPublished bool
	CancelSignalSent bool
	CleanupScheduled bool
}

var (
	exportJobs      = map[uint32]*ExportJob{}
	exportJobsMu    sync.RWMutex
	nextExportJobID atomic.Uint32
)

var errExportJobCancelled = errors.New("export cancelled")

func StartTimeSeriesExportAggregateKPI(userID uint32, input graphQLModel.TimeSeriesReadAggregateKPIInput) (graphQLModel.TimeSeriesExport, error) {
	req, err := buildTimeSeriesAggregateReadRequest(userID, input)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return StartTimeSeriesExportBase(userID, req, input.SdTypeID)
}

func StartTimeSeriesExport(userID uint32, input graphQLModel.TimeSeriesReadInput) (graphQLModel.TimeSeriesExport, error) {
	req, err := buildTimeSeriesReadRequest(userID, input)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return StartTimeSeriesExportBase(userID, req, input.SdTypeID)
}

func StartTimeSeriesExportBase(userID uint32, req sharedModel.TimeSeriesReadRequest, sdTypeID *uint32) (graphQLModel.TimeSeriesExport, error) {
	id := nextExportJobID.Add(1)
	createdAt := time.Now().UTC()
	req.JobID = id
	job := &ExportJob{
		ID:        id,
		UserID:    userID,
		FilePath:  fmt.Sprintf("/tmp/export_%d.csv", id),
		Done:      make(chan struct{}),
		Status:    graphQLModel.ExportStatusPending,
		CreatedAt: createdAt,
	}
	exportJobsMu.Lock()
	exportJobs[id] = job
	exportJobsMu.Unlock()
	publishExportJobUpdate(job)
	go runExportJob(job, req, sdTypeID)
	return snapshotExportJob(job), nil
}

func GetTimeSeriesExport(userID uint32, id uint32) (graphQLModel.TimeSeriesExport, error) {
	job, err := getOwnedExportJob(userID, id)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	return snapshotExportJob(job), nil
}

func CancelTimeSeriesExport(userID uint32, id uint32) (graphQLModel.TimeSeriesExport, error) {
	job, err := getOwnedExportJob(userID, id)
	if err != nil {
		return graphQLModel.TimeSeriesExport{}, err
	}
	exportJobsMu.Lock()
	if isTerminalExportStatus(job.Status) {
		snapshot := snapshotExportJobLocked(job)
		exportJobsMu.Unlock()
		return snapshot, nil
	}
	job.CancelRequested = true
	job.Status = graphQLModel.ExportStatusCancelled
	job.Err = errExportJobCancelled
	job.ErrorMessage = nil
	setJobExpiryLocked(job)
	shouldSignalCancel := job.RequestPublished && !job.CancelSignalSent
	if shouldSignalCancel {
		job.CancelSignalSent = true
	}
	snapshot := snapshotExportJobLocked(job)
	exportJobsMu.Unlock()
	if shouldSignalCancel {
		logExportJobCancelSignalFailure(id, publishTimeSeriesReadCancel(id))
	}
	publishExportJobUpdate(job)
	scheduleCleanup(job)
	return snapshot, nil
}

func runExportJob(job *ExportJob, req sharedModel.TimeSeriesReadRequest, sdTypeID *uint32) {
	defer close(job.Done)
	req.Limit = nil
	req.Batch = nil
	if isExportJobCancelled(job) {
		return
	}
	setExportJobState(job, graphQLModel.ExportStatusProcessing, nil, nil)
	session, err := newTimeSeriesRPCSession(fmt.Sprint(job.ID))
	if err != nil {
		failExportJob(job, err)
		return
	}
	defer session.Close()
	handler, closeHandler := buildExportHandler(job, req, sdTypeID)
	if handler == nil {
		return
	}
	defer closeHandler()
	if err := session.PublishReadRequest(req); err != nil {
		failExportJob(job, err)
		return
	}
	markExportJobRequestPublished(job)
	err = rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesReadResponse](session.messages, session.correlationID, 0, func(resp sharedModel.TimeSeriesReadResponse, msg amqp091.Delivery) (bool, error) {
		if isExportJobCancelled(job) {
			return true, nil
		}
		if resp.Error != "" {
			err := fmt.Errorf("%s", resp.Error)
			failExportJob(job, err)
			return true, err
		}
		if err := handler(resp, msg); err != nil {
			return true, err
		}
		if !resp.HasMoreBatches {
			downloadURL := fmt.Sprintf("/rest/time-series/export/%d", job.ID)
			setExportJobState(job, graphQLModel.ExportStatusDone, &downloadURL, nil)
			scheduleCleanup(job)
			return true, nil
		}
		return false, nil
	})
	if err != nil && !errors.Is(err, errExportJobCancelled) {
		failExportJob(job, err)
	}
}

func buildExportHandler(job *ExportJob, req sharedModel.TimeSeriesReadRequest, sdTypeID *uint32) (func(resp sharedModel.TimeSeriesReadResponse, _ amqp091.Delivery) error, func()) {
	file, err := os.Create(job.FilePath)
	if err != nil {
		failExportJob(job, err)
		return nil, func() {}
	}
	writer := csv.NewWriter(file)
	params := loadParametersFromDB(graphQLModel.TimeSeriesType(req.Type), sdTypeID)
	headerWritten := false
	closed := false
	closeResources := func() {
		if closed {
			return
		}
		writer.Flush()
		_ = file.Close()
		closed = true
	}
	return func(resp sharedModel.TimeSeriesReadResponse, _ amqp091.Delivery) error {
		defer func() {
			if !resp.HasMoreBatches {
				closeResources()
			}
		}()
		if isExportJobCancelled(job) {
			return errExportJobCancelled
		}
		if resp.Error != "" {
			err := fmt.Errorf("%s", resp.Error)
			failExportJob(job, err)
			return err
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
	}, closeResources
}

func GetExportJob(id uint32) (*ExportJob, bool) {
	exportJobsMu.RLock()
	defer exportJobsMu.RUnlock()
	job, ok := exportJobs[id]
	return job, ok
}

func publishExportJobUpdate(job *ExportJob) {
	events.GetEventBus().Publish(events.TimeSeriesExportUpdatedEventType, snapshotExportJob(job))
}

func snapshotExportJob(job *ExportJob) graphQLModel.TimeSeriesExport {
	exportJobsMu.RLock()
	defer exportJobsMu.RUnlock()
	return snapshotExportJobLocked(job)
}

func snapshotExportJobLocked(job *ExportJob) graphQLModel.TimeSeriesExport {
	var expiresAt *string
	if job.ExpiresAt != nil {
		value := job.ExpiresAt.UTC().Format(time.RFC3339Nano)
		expiresAt = &value
	}
	return graphQLModel.TimeSeriesExport{
		ID:          job.ID,
		Status:      job.Status,
		DownloadURL: job.DownloadURL,
		CreatedAt:   job.CreatedAt.UTC().Format(time.RFC3339Nano),
		ExpiresAt:   expiresAt,
		Error:       job.ErrorMessage,
	}
}

func setExportJobState(job *ExportJob, status graphQLModel.ExportStatus, downloadURL *string, err error) {
	exportJobsMu.Lock()
	job.Status = status
	job.DownloadURL = downloadURL
	if err != nil {
		message := err.Error()
		job.Err = err
		job.ErrorMessage = &message
	} else {
		job.Err = nil
		job.ErrorMessage = nil
	}
	if isTerminalExportStatus(status) {
		setJobExpiryLocked(job)
	}
	exportJobsMu.Unlock()
	publishExportJobUpdate(job)
}

func failExportJob(job *ExportJob, err error) {
	if err == nil || errors.Is(err, errExportJobCancelled) || isExportJobCancelled(job) {
		return
	}
	setExportJobState(job, graphQLModel.ExportStatusFailed, nil, err)
	scheduleCleanup(job)
}

func getOwnedExportJob(userID uint32, id uint32) (*ExportJob, error) {
	exportJobsMu.RLock()
	defer exportJobsMu.RUnlock()
	job, ok := exportJobs[id]
	if !ok {
		return nil, fmt.Errorf("export not found")
	}
	if job.UserID != userID {
		return nil, fmt.Errorf("export not found")
	}
	return job, nil
}

func isExportJobCancelled(job *ExportJob) bool {
	exportJobsMu.RLock()
	defer exportJobsMu.RUnlock()
	return job.CancelRequested || job.Status == graphQLModel.ExportStatusCancelled
}

func isTerminalExportStatus(status graphQLModel.ExportStatus) bool {
	switch status {
	case graphQLModel.ExportStatusDone, graphQLModel.ExportStatusFailed, graphQLModel.ExportStatusCancelled, graphQLModel.ExportStatusExpired:
		return true
	default:
		return false
	}
}

func setJobExpiryLocked(job *ExportJob) {
	expiresAt := time.Now().UTC().Add(30 * time.Minute)
	job.ExpiresAt = &expiresAt
}

func scheduleCleanup(job *ExportJob) {
	exportJobsMu.Lock()
	if job.CleanupScheduled {
		exportJobsMu.Unlock()
		return
	}
	job.CleanupScheduled = true
	exportJobsMu.Unlock()
	go cleanupJob(job.ID)
}

func cleanupJob(id uint32) {
	time.Sleep(30 * time.Minute)
	exportJobsMu.Lock()
	job := exportJobs[id]
	delete(exportJobs, id)
	exportJobsMu.Unlock()
	if job != nil {
		_ = os.Remove(job.FilePath)
	}
}

func markExportJobRequestPublished(job *ExportJob) {
	exportJobsMu.Lock()
	job.RequestPublished = true
	shouldSignalCancel := job.CancelRequested && !job.CancelSignalSent
	if shouldSignalCancel {
		job.CancelSignalSent = true
	}
	exportJobsMu.Unlock()
	if shouldSignalCancel {
		logExportJobCancelSignalFailure(job.ID, publishTimeSeriesReadCancel(job.ID))
	}
}

func publishTimeSeriesReadCancel(jobID uint32) error {
	payload, err := json.Marshal(sharedModel.TimeSeriesReadCancelRequest{
		JobID: jobID,
	})
	if err != nil {
		return err
	}
	return getDLLRabbitMQClient().PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesReadCancelRequestQueueName), payload)
}

func logExportJobCancelSignalFailure(jobID uint32, err error) {
	if err != nil {
		log.Printf("[EXPORT] failed to publish remote cancel | jobID=%d err=%s", jobID, err.Error())
	}
}
