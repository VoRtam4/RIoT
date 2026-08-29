/**
 * @file reprocess.go
 * @brief Řízení přepočtu KPI nad historickými surovými daty.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_message_processing_unit
 */
package processing

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ReprocessKPI(req sharedModel.KPIReprocessRequestISCMessage) error {
	client := rabbitmq.NewClient()
	defer client.Dispose()
	waitForConfig(req.JobID)
	key := makeKey(req.SDTypeUID, req.KPIDefinitionUID)
	activeReprocessJobs.Store(key, req.JobID)
	kpiDefinitionsBySDTypeDenotationMapMutex.RLock()
	kpis := kpiDefinitionsBySDTypeDenotationMap[req.SDTypeUID]
	kpiDefinitionsBySDTypeDenotationMapMutex.RUnlock()
	var targetKPI *sharedModel.KPIDefinitionMPU
	for i := range kpis {
		if kpis[i].UID == req.KPIDefinitionUID {
			targetKPI = &kpis[i]
			break
		}
	}
	if targetKPI == nil {
		return fmt.Errorf("KPI definition not found")
	}
	readReq := sharedModel.TimeSeriesReprocessReadRequest{
		JobID:            req.JobID,
		Wait:             req.Wait,
		SDTypeUID:        req.SDTypeUID,
		SDInstanceUIDs:   req.SDInstanceUIDs,
		KPIDefinitionUID: req.KPIDefinitionUID,
		To:               req.To,
		Batch:            limit,
	}
	jsonReq := sharedUtils.SerializeToJSON(readReq)
	if jsonReq.IsFailure() {
		return jsonReq.GetError()
	}
	ch := client.GetChannel()
	replyQueue, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(replyQueue.Name, "", false, true, false, false, nil)
	if err != nil {
		return err
	}
	correlationID := uuid.New().String()
	err = ch.PublishWithContext(context.Background(), "", sharedConstants.TimeSeriesReprocessReadRequestQueueName, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Body:          jsonReq.GetPayload(),
		CorrelationId: correlationID,
		ReplyTo:       replyQueue.Name,
	})
	if err != nil {
		return err
	}
	instanceState := map[string]sharedModel.KPIState{}
	rawInstanceState := map[string]map[string]interface{}{}
	handler := createReprocessHandler(client, req, key, targetKPI, instanceState, rawInstanceState)
	return rabbitmq.ConsumeRPCStream[sharedModel.TimeSeriesReprocessReadResponse](msgs, correlationID, 0, func(resp sharedModel.TimeSeriesReprocessReadResponse, msg amqp091.Delivery) (bool, error) {
		if resp.Error != "" {
			return true, fmt.Errorf("%s", resp.Error)
		}
		if err := handler(resp, msg); err != nil {
			return true, err
		}
		if !resp.HasMore {
			return true, nil
		}
		return false, nil
	})
}

func createReprocessHandler(client rabbitmq.Client, req sharedModel.KPIReprocessRequestISCMessage, key string, targetKPI *sharedModel.KPIDefinitionMPU, instanceState map[string]sharedModel.KPIState, rawInstanceState map[string]map[string]interface{}) func(resp sharedModel.TimeSeriesReprocessReadResponse, delivery amqp.Delivery) error {
	return func(resp sharedModel.TimeSeriesReprocessReadResponse, delivery amqp.Delivery) error {
		if !isActive(key, req.JobID) {
			log.Printf("[MPU][REPROCESS] Job cancelled before batch processing | kpiUID=%s jobID=%s", req.KPIDefinitionUID, req.JobID)
			return nil
		}
		if resp.Error != "" {
			return fmt.Errorf("%s", resp.Error)
		}
		tsBatch := make([]sharedModel.TimeSeriesKPIResultRecord, 0, len(resp.Data))
		for _, point := range resp.Data {
			if !isActive(key, req.JobID) {
				log.Printf("[MPU][REPROCESS] Job cancelled during point processing | kpiUID=%s jobID=%s", req.KPIDefinitionUID, req.JobID)
				return nil
			}
			sdInstanceUID := point.Tags["sdInstanceUID"]
			previousRaw := rawInstanceState[sdInstanceUID]
			result, currentRaw, ok := evaluateReprocessPoint(point, req.KPIDefinitionUID, targetKPI, previousRaw)
			if !ok {
				continue
			}
			rawInstanceState[result.SDInstanceUID] = currentRaw
			prev, exists := instanceState[result.SDInstanceUID]
			if exists && prev.Value == result.Fulfilled {
				prev.SynchronizedAt = time.Now()
				instanceState[result.SDInstanceUID] = prev
				continue
			}
			instanceState[result.SDInstanceUID] = sharedModel.KPIState{
				Value:          result.Fulfilled,
				EventTime:      result.EventTime,
				SynchronizedAt: time.Now(),
			}
			tags := make(map[string]string)
			for k, v := range point.Tags {
				switch k {
				case "sdInstanceUID", "sdType", "kpiDefinitionUID":
					continue
				default:
					tags[k] = v
				}
			}
			tsBatch = append(tsBatch, sharedModel.TimeSeriesKPIResultRecord{
				JobID:            req.JobID,
				EventTime:        result.EventTime,
				SDInstanceUID:    result.SDInstanceUID,
				SDTypeUID:        req.SDTypeUID,
				KPIDefinitionUID: result.KPIDefinitionUID,
				Fulfilled:        result.Fulfilled,
				Tags:             tags,
			})
		}
		if len(tsBatch) > 0 {
			if !isActive(key, req.JobID) {
				log.Printf("[MPU][REPROCESS] Job cancelled before write | kpiUID=%s jobID=%s", req.KPIDefinitionUID, req.JobID)
				return nil
			}
			storeKPI(client, tsBatch)
		}
		if !resp.HasMore {
			if req.JobID != "" && !isActive(key, req.JobID) {
				log.Printf("[MPU][REPROCESS] Job cancelled before finalize | kpiUID=%s jobID=%s", req.KPIDefinitionUID, req.JobID)
				return nil
			}
			finalizeReprocess(client, req, instanceState)
			return nil
		}
		return nil
	}
}

func evaluateReprocessPoint(point sharedModel.TimeSeriesDataPoint, kpiUID string, targetKPI *sharedModel.KPIDefinitionMPU, previousRaw map[string]interface{}) (sharedModel.KPIFulfillmentCheckResultISCMessage, map[string]interface{}, bool) {
	sdInstanceUID := point.Tags["sdInstanceUID"]
	if sdInstanceUID == "" {
		log.Printf("[MPU][REPROCESS] Missing sdInstanceUID -> skipping point")
		return sharedModel.KPIFulfillmentCheckResultISCMessage{}, nil, false
	}
	if targetKPI.SDInstanceMode == sharedModel.SELECTED {
		allowed := false
		for _, uid := range targetKPI.SelectedSDInstanceUIDs {
			if uid == sdInstanceUID {
				allowed = true
				break
			}
		}
		if !allowed {
			log.Printf("[MPU][REPROCESS] Instance not allowed for KPI | uid=%s", sdInstanceUID)
			return sharedModel.KPIFulfillmentCheckResultISCMessage{}, nil, false
		}
	}
	currentRaw := mergeTagsAndFields(point)
	if previousRaw == nil {
		previousRaw = map[string]interface{}{}
	}
	result := CheckKPIFulfillment(*targetKPI, KPIEvaluationContext{
		CurrentValues:  currentRaw,
		PreviousValues: previousRaw,
	})
	if result.IsFailure() {
		log.Printf("[MPU][REPROCESS] KPI evaluation failed: %v", result.GetError())
		return sharedModel.KPIFulfillmentCheckResultISCMessage{}, nil, false
	}
	return sharedModel.KPIFulfillmentCheckResultISCMessage{
		EventTime:        point.Time,
		SDInstanceUID:    sdInstanceUID,
		KPIDefinitionUID: kpiUID,
		Fulfilled:        result.GetPayload(),
	}, currentRaw, true
}

func finalizeReprocess(client rabbitmq.Client, req sharedModel.KPIReprocessRequestISCMessage, instanceState map[string]sharedModel.KPIState) {
	results := make([]sharedModel.KPIFulfillmentCheckResultISCMessage, 0, limit)
	for uid, state := range instanceState {
		key := sharedModel.KPIKey{
			SDInstanceUID:    uid,
			KPIDefinitionUID: req.KPIDefinitionUID,
		}
		lastStateRaw, exists := lastKPI.Load(key)
		if exists {
			lastState, ok := lastStateRaw.(sharedModel.KPIState)
			if ok && lastState.SynchronizedAt.After(req.To) {
				log.Printf("[MPU][REPROCESS] Skipping instance in finalize, newer event exists | uid=%s", uid)
				continue
			}
		}
		lastKPI.Store(key, state)
		results = append(results, sharedModel.KPIFulfillmentCheckResultISCMessage{
			SDTypeUID:        req.SDTypeUID,
			EventTime:        state.EventTime,
			SDInstanceUID:    uid,
			KPIDefinitionUID: req.KPIDefinitionUID,
			Fulfilled:        state.Value,
		})
		if len(results) >= limit {
			publishKPI(client, results, true)
			results = results[:0]
		}
	}
	if len(results) > 0 {
		publishKPI(client, results, true)
	}
}
