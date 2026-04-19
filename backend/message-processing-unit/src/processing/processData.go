package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	kpiDefinitionsBySDTypeDenotationMap      map[string][]sharedModel.KPIDefinitionMPU = make(map[string][]sharedModel.KPIDefinitionMPU)
	kpiDefinitionsBySDTypeDenotationMapMutex sync.RWMutex
	unitUUID                                 string
	limit                                    = 500
	sdTypeDefinitions                        = make(map[string]sharedModel.SDTypeDefinitionCache)
	sdTypeDefinitionsMutex                   sync.RWMutex
	activeReprocessJobs                      sync.Map
	jobSyncMap                               sync.Map
)

func InitializeProcessing() {
	unitUUID = uuid.New().String()
	log.Printf("[MPU] Processing unit initialized | unitUUID=%s", unitUUID)
}

func DenotationMapUpdates(rabbitMQClient rabbitmq.Client) {
	log.Printf("[MPU] Waiting for KPI configuration updates | unit=%s", unitUUID)
	err := rabbitmq.ConsumeJSONMessagesFromFanoutExchange[sharedModel.KPIConfigurationUpdateISCMessage](rabbitMQClient, sharedConstants.BuiltInFanoutExchangeName, func(messagePayload sharedModel.KPIConfigurationUpdateISCMessage) error {
		log.Printf("[MPU] KPI configuration update received | unit=%s entries=%d", unitUUID, len(messagePayload.KpiConfiguration))
		kpiDefinitionsBySDTypeDenotationMapMutex.Lock()
		kpiDefinitionsBySDTypeDenotationMap = messagePayload.KpiConfiguration
		kpiDefinitionsBySDTypeDenotationMapMutex.Unlock()
		log.Printf("[MPU] KPI configuration map updated successfully | unit=%s", unitUUID)
		signalConfigReady(messagePayload.JobID)
		return nil
	})
	if err != nil {
		log.Printf("[MPU] Consumption of messages from the '%s' fanout exchange has failed: %s\n", sharedConstants.BuiltInFanoutExchangeName, err.Error())
	}
}

func UpdateSDType(messages []sharedModel.SDTypeUpdateISCMessage) {
	log.Printf("[MPU][SDTYPE] Received SDType update | count=%d", len(messages))
	newDefinitions := make(map[string]sharedModel.SDTypeDefinitionCache)
	for _, msg := range messages {
		paramMap := make(map[string]sharedModel.SDParameter, len(msg.Parameters))
		for _, p := range msg.Parameters {
			paramMap[p.Denotation] = sharedModel.SDParameter{
				Denotation: p.Denotation,
				Type:       sharedModel.SDParameterType(p.Type),
				Label:      p.Label,
				Role:       sharedModel.SDParameterRole(p.Role),
			}
		}
		newDefinitions[msg.SDTypeUID] = sharedModel.SDTypeDefinitionCache{
			Label:      msg.SDTypeUID,
			Parameters: paramMap,
		}
	}
	sdTypeDefinitionsMutex.Lock()
	sdTypeDefinitions = newDefinitions
	sdTypeDefinitionsMutex.Unlock()
	log.Printf("[MPU][SDTYPE] SDType cache updated | types=%d", len(newDefinitions))
}

func ProcessRaw(rabbitMQClient rabbitmq.Client, messagePayload sharedModel.KPIFulfillmentCheckRequestISCMessage, params map[string]interface{}, eventTime time.Time) bool {
	uid := messagePayload.SDInstanceUID
	lastSnapshotRaw, exists := lastRaw.Load(uid)
	var lastSnapshot sharedModel.RawState
	if exists {
		lastSnapshot = lastSnapshotRaw.(sharedModel.RawState)
	} else {
		lastSnapshot = sharedModel.RawState{
			Values: map[string]interface{}{},
		}
	}
	result := processParams(messagePayload.SDTypeUID, params, lastSnapshot.Values)
	if result.SDTypeChanged {
		publishSDTypeUpdate(rabbitMQClient, messagePayload.SDTypeUID)
	}
	newValues := result.Values
	changed := result.Changed
	if exists {
		lastSnapshot.SynchronizedAt = time.Now()
		if !changed {
			lastRaw.Store(uid, lastSnapshot)
		}
	}
	fields, tags := splitParamsBySDType(messagePayload.SDTypeUID, newValues)
	if exists && eventTime.Before(lastSnapshot.EventTime) {
		log.Printf("[MPU][RAW] Older event -> storing ONLY TS | uid=%s", uid)
		toStore := []sharedModel.TimeSeriesRawRecord{{
			EventTime:     eventTime,
			SDInstanceUID: uid,
			SDTypeUID:     messagePayload.SDTypeUID,
			Fields:        fields,
			Tags:          tags,
		}}
		storeRaw(rabbitMQClient, toStore)
		return true
	}
	if !changed {
		log.Printf("[MPU][RAW] No change detected -> skipping RAW write | uid=%s", uid)
		return false
	}
	payloadResult := sharedUtils.SerializeToJSON(newValues)
	if payloadResult.IsFailure() {
		log.Printf("[MPU][RAW] Payload serialization error: %s", payloadResult.GetError())
		return false
	}
	publishRaw(rabbitMQClient, []sharedModel.RawDataPointISCMessage{{
		SDTypeUID:     messagePayload.SDTypeUID,
		SDInstanceUID: uid,
		EventTime:     eventTime,
		Payload:       payloadResult.GetPayload(),
	}})
	lastRaw.Store(uid, sharedModel.RawState{
		Values:         newValues,
		EventTime:      eventTime,
		SynchronizedAt: time.Now(),
	})
	toStore := []sharedModel.TimeSeriesRawRecord{{
		EventTime:     eventTime,
		SDInstanceUID: uid,
		SDTypeUID:     messagePayload.SDTypeUID,
		Fields:        fields,
		Tags:          tags,
	}}
	storeRaw(rabbitMQClient, toStore)
	return true
}

func processParams(sdTypeUID string, params map[string]interface{}, lastValues map[string]interface{}) sharedModel.ProcessResult {
	sdTypeDefinitionsMutex.RLock()
	definitionItem, exists := sdTypeDefinitions[sdTypeUID]
	sdTypeDefinitionsMutex.RUnlock()
	newValues := make(map[string]interface{})
	changed := false
	sdTypeChanged := false
	keys := make(map[string]struct{})
	for k := range params {
		keys[k] = struct{}{}
	}
	for k := range lastValues {
		keys[k] = struct{}{}
	}
	for k := range keys {
		v, existsNow := params[k]
		var normVal interface{}
		var detectedType sharedModel.SDParameterType
		var arrived bool
		if existsNow && v != nil {
			var ok bool
			normVal, detectedType, ok = normalizeValueWithType(v)
			if ok {
				arrived = true
			}
		}
		var paramDef sharedModel.SDParameter
		var paramExists bool
		if exists {
			paramDef, paramExists = definitionItem.Parameters[k]
		}
		if arrived && (!exists || !paramExists) {
			sdTypeDefinitionsMutex.Lock()
			definitionItem, exists = sdTypeDefinitions[sdTypeUID]
			if !exists {
				definitionItem = sharedModel.SDTypeDefinitionCache{
					Label:      sdTypeUID,
					Parameters: make(map[string]sharedModel.SDParameter),
				}
			}
			paramDef, paramExists = definitionItem.Parameters[k]
			if !paramExists {
				paramDef = sharedModel.SDParameter{
					Denotation: k,
					Type:       detectedType,
					Label:      sharedUtils.SafeLabel("", k),
					Role:       sharedModel.SDParameterRoleField,
				}
				definitionItem.Parameters[k] = paramDef
				sdTypeDefinitions[sdTypeUID] = definitionItem
				sdTypeChanged = true
			}
			sdTypeDefinitionsMutex.Unlock()
		}
		if arrived && paramExists && paramDef.Type != detectedType {
			arrived = false
		}
		oldVal, existedBefore := lastValues[k]
		if !existedBefore && !arrived {
			continue
		}
		if existedBefore && !arrived {
			if oldVal != nil {
				changed = true
			}
			newValues[k] = nil
			continue
		}
		if !existedBefore && arrived {
			changed = true
			newValues[k] = normVal
			continue
		}
		if !sharedUtils.CompareJSONs(oldVal, normVal) {
			changed = true
		}
		newValues[k] = normVal
	}
	return sharedModel.ProcessResult{
		Values:        newValues,
		Changed:       changed,
		SDTypeChanged: sdTypeChanged,
	}
}

func normalizeValueWithType(v interface{}) (interface{}, sharedModel.SDParameterType, bool) {

	switch val := v.(type) {

	case string:
		return val, sharedModel.SDParameterTypeString, true

	case bool:
		return val, sharedModel.SDParameterTypeBoolean, true

	case int:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case int8:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case int16:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case int32:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case int64:
		return float64(val), sharedModel.SDParameterTypeNumber, true

	case uint:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case uint8:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case uint16:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case uint32:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case uint64:
		return float64(val), sharedModel.SDParameterTypeNumber, true

	case float32:
		return float64(val), sharedModel.SDParameterTypeNumber, true
	case float64:
		return val, sharedModel.SDParameterTypeNumber, true

	case nil:
		return nil, sharedModel.SDParameterTypeString, false

	default:
		b, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprint(val), sharedModel.SDParameterTypeString, true
		}
		return string(b), sharedModel.SDParameterTypeString, true
	}
}

func ProcessKPI(rabbitMQClient rabbitmq.Client, messagePayload sharedModel.KPIFulfillmentCheckRequestISCMessage, params map[string]interface{}, eventTime time.Time) error {
	kpiDefinitionsBySDTypeDenotationMapMutex.RLock()
	kpiDefinitions := kpiDefinitionsBySDTypeDenotationMap[messagePayload.SDTypeUID]
	kpiDefinitionsBySDTypeDenotationMapMutex.RUnlock()
	sdInstanceUID := messagePayload.SDInstanceUID
	results := sharedUtils.EmptySlice[sharedModel.KPIFulfillmentCheckResultISCMessage]()
	tsRecords := sharedUtils.EmptySlice[sharedModel.TimeSeriesKPIResultRecord]()
	_, tags := splitParamsBySDType(messagePayload.SDTypeUID, params)
	var paramAny any = params
	for _, kpiDefinition := range kpiDefinitions {
		if kpiDefinition.SDInstanceMode != sharedModel.ALL {
			containsUID := false
			for _, uid := range kpiDefinition.SelectedSDInstanceUIDs {
				if uid == sdInstanceUID {
					containsUID = true
					break
				}
			}
			if !containsUID {
				continue
			}
		}
		result := CheckKPIFulfillment(kpiDefinition, &paramAny)
		if result.IsFailure() {
			continue
		}
		value := result.GetPayload()
		kpiID := sharedUtils.NewOptionalFromPointer(kpiDefinition.ID).GetPayload()
		key := sharedModel.KPIKey{
			SDInstanceUID:   sdInstanceUID,
			KPIDefinitionID: kpiID,
		}
		last, exists := lastKPI.Load(key)
		if exists {
			state := last.(sharedModel.KPIState)
			if !eventTime.Before(state.EventTime) && state.Value == value {
				state.SynchronizedAt = time.Now()
				lastKPI.Store(key, state)
				log.Printf("[MPU][KPI] KPI unchanged -> skipping store | kpiID=%d uid=%s", kpiID, sdInstanceUID)
				continue
			}
		}
		lastKPI.Store(key, sharedModel.KPIState{
			Value:          value,
			EventTime:      eventTime,
			SynchronizedAt: time.Now(),
		})
		results = append(results, sharedModel.KPIFulfillmentCheckResultISCMessage{
			SDTypeUID:       messagePayload.SDTypeUID,
			EventTime:       eventTime,
			SDInstanceUID:   sdInstanceUID,
			KPIDefinitionID: kpiID,
			Fulfilled:       value,
		})
		toStore := sharedModel.TimeSeriesKPIResultRecord{
			JobID:           "",
			EventTime:       eventTime,
			SDInstanceUID:   sdInstanceUID,
			SDTypeUID:       messagePayload.SDTypeUID,
			KPIDefinitionID: kpiID,
			Fulfilled:       value,
			Tags:            tags,
		}
		tsRecords = append(tsRecords, toStore)
	}
	storeKPI(rabbitMQClient, tsRecords)
	publishKPI(rabbitMQClient, results, false)
	return nil
}

func ReprocessKPI(req sharedModel.KPIReprocessRequestISCMessage) error {
	client := rabbitmq.NewClient()
	defer client.Dispose()
	waitForConfig(req.JobID)
	key := makeKey(req.SDTypeUID, req.KPIDefinitionID)
	activeReprocessJobs.Store(key, req.JobID)
	kpiDefinitionsBySDTypeDenotationMapMutex.RLock()
	kpis := kpiDefinitionsBySDTypeDenotationMap[req.SDTypeUID]
	kpiDefinitionsBySDTypeDenotationMapMutex.RUnlock()
	var targetKPI *sharedModel.KPIDefinitionMPU
	for i := range kpis {
		if kpis[i].ID != nil && *kpis[i].ID == req.KPIDefinitionID {
			targetKPI = &kpis[i]
			break
		}
	}
	if targetKPI == nil {
		return fmt.Errorf("KPI definition not found")
	}
	readReq := sharedModel.TimeSeriesReprocessReadRequest{
		JobID:           req.JobID,
		Wait:            req.Wait,
		SDTypeUID:       req.SDTypeUID,
		SDInstanceUIDs:  req.SDInstanceUIDs,
		KPIDefinitionID: req.KPIDefinitionID,
		To:              req.To,
		Batch:           limit,
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
	handler := createReprocessHandler(client, req, key, targetKPI, instanceState)
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

func createReprocessHandler(client rabbitmq.Client, req sharedModel.KPIReprocessRequestISCMessage, key string, targetKPI *sharedModel.KPIDefinitionMPU, instanceState map[string]sharedModel.KPIState) func(resp sharedModel.TimeSeriesReprocessReadResponse, delivery amqp.Delivery) error {
	return func(resp sharedModel.TimeSeriesReprocessReadResponse, delivery amqp.Delivery) error {
		if !isActive(key, req.JobID) {
			log.Printf("[MPU][REPROCESS] Job cancelled before batch processing | kpiID=%d jobID=%s", req.KPIDefinitionID, req.JobID)
			return nil
		}
		if resp.Error != "" {
			return fmt.Errorf("%s", resp.Error)
		}
		tsBatch := make([]sharedModel.TimeSeriesKPIResultRecord, 0, len(resp.Data))
		for _, point := range resp.Data {
			if !isActive(key, req.JobID) {
				log.Printf("[MPU][REPROCESS] Job cancelled during point processing | kpiID=%d jobID=%s", req.KPIDefinitionID, req.JobID)
				return nil
			}
			result, ok := evaluateReprocessPoint(point, req.KPIDefinitionID, targetKPI)
			if !ok {
				continue
			}
			prev, exists := instanceState[result.SDInstanceUID]
			if exists && prev.Value == result.Fulfilled {
				prev.SynchronizedAt = time.Now()
				instanceState[result.SDInstanceUID] = prev
				log.Printf("[MPU][REPROCESS] Skipping point, duplicated result | kpiID=%d jobID=%s", req.KPIDefinitionID, req.JobID)
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
				case "sdInstanceUID", "sdType", "kpiDefinitionID":
					continue
				default:
					tags[k] = v
				}
			}
			tsBatch = append(tsBatch, sharedModel.TimeSeriesKPIResultRecord{
				JobID:           req.JobID,
				EventTime:       result.EventTime,
				SDInstanceUID:   result.SDInstanceUID,
				SDTypeUID:       req.SDTypeUID,
				KPIDefinitionID: result.KPIDefinitionID,
				Fulfilled:       result.Fulfilled,
				Tags:            tags,
			})
		}
		if len(tsBatch) > 0 {
			if !isActive(key, req.JobID) {
				log.Printf("[MPU][REPROCESS] Job cancelled before write | kpiID=%d jobID=%s", req.KPIDefinitionID, req.JobID)
				return nil
			}
			storeKPI(client, tsBatch)
		}
		if !resp.HasMore {
			if req.JobID != "" && !isActive(key, req.JobID) {
				log.Printf("[MPU][REPROCESS] Job cancelled before finalize | kpiID=%d jobID=%s", req.KPIDefinitionID, req.JobID)
				return nil
			}
			finalizeReprocess(client, req, instanceState)
			return nil
		}
		return nil
	}
}

func evaluateReprocessPoint(point sharedModel.TimeSeriesDataPoint, kpiID uint32, targetKPI *sharedModel.KPIDefinitionMPU) (sharedModel.KPIFulfillmentCheckResultISCMessage, bool) {
	sdInstanceUID := point.Tags["sdInstanceUID"]
	if sdInstanceUID == "" {
		log.Printf("[MPU][REPROCESS] Missing sdInstanceUID -> skipping point")
		return sharedModel.KPIFulfillmentCheckResultISCMessage{}, false
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
			return sharedModel.KPIFulfillmentCheckResultISCMessage{}, false
		}
	}
	var paramAny any = mergeTagsAndFields(point)
	result := CheckKPIFulfillment(*targetKPI, &paramAny)
	if result.IsFailure() {
		log.Printf("[MPU][REPROCESS] KPI evaluation failed: %v", result.GetError())
		return sharedModel.KPIFulfillmentCheckResultISCMessage{}, false
	}
	return sharedModel.KPIFulfillmentCheckResultISCMessage{
		EventTime:       point.Time,
		SDInstanceUID:   sdInstanceUID,
		KPIDefinitionID: kpiID,
		Fulfilled:       result.GetPayload(),
	}, true
}

func finalizeReprocess(client rabbitmq.Client, req sharedModel.KPIReprocessRequestISCMessage, instanceState map[string]sharedModel.KPIState) {
	results := make([]sharedModel.KPIFulfillmentCheckResultISCMessage, 0, limit)
	for uid, state := range instanceState {
		key := sharedModel.KPIKey{
			SDInstanceUID:   uid,
			KPIDefinitionID: req.KPIDefinitionID,
		}
		lastKPI.Store(key, state)
		results = append(results, sharedModel.KPIFulfillmentCheckResultISCMessage{
			SDTypeUID:       req.SDTypeUID,
			EventTime:       state.EventTime,
			SDInstanceUID:   uid,
			KPIDefinitionID: req.KPIDefinitionID,
			Fulfilled:       state.Value,
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

func publishSDTypeUpdate(client rabbitmq.Client, sdTypeUID string) {
	sdTypeDefinitionsMutex.RLock()
	definition, exists := sdTypeDefinitions[sdTypeUID]
	sdTypeDefinitionsMutex.RUnlock()
	if !exists {
		log.Printf("[MPU][SDTYPE] Definition not found for type=%s", sdTypeUID)
		return
	}
	params := make([]sharedModel.SDParameter, 0, len(definition.Parameters))
	for _, p := range definition.Parameters {
		params = append(params, p)
	}
	label := definition.Label
	if label == "" {
		label = strings.ToUpper(sdTypeUID[:1]) + sdTypeUID[1:]
	}
	msg := sharedModel.SDTypeRegistrationRequestISCMessage{
		SDTypeUID:  sdTypeUID,
		Label:      label,
		Parameters: params,
	}
	jsonResult := sharedUtils.SerializeToJSON(msg)
	if jsonResult.IsFailure() {
		log.Printf("[MPU][SDTYPE] Serialization error: %s", jsonResult.GetError())
		return
	}
	err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.SDTypeRegistrationRequestsQueueName), jsonResult.GetPayload())
	if err != nil {
		log.Printf("[MPU][SDTYPE] Failed publishing SDType registration/upsert: %s", err)
		return
	}
	log.Printf("[MPU][SDTYPE] SDType registration/upsert published | type=%s label=%s params=%d", sdTypeUID, label, len(params))
}

func publishRaw(client rabbitmq.Client, results []sharedModel.RawDataPointISCMessage) {
	if len(results) == 0 {
		return
	}
	jsonResult := sharedUtils.SerializeToJSON(results)
	if jsonResult.IsFailure() {
		log.Printf("[MPU][RAW] Serialization error: %s", jsonResult.GetError())
		return
	}
	err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.RawDataPointQueueName), jsonResult.GetPayload())
	if err != nil {
		log.Printf("[MPU][RAW] Failed publishing RAW data: %s", err)
		return
	}
}

func publishKPI(client rabbitmq.Client, results []sharedModel.KPIFulfillmentCheckResultISCMessage, reprocess bool) {
	if len(results) == 0 {
		return
	}
	jsonResult := sharedUtils.SerializeToJSON(sharedModel.KPIFulfillmentCheckResultTupleISCMessage{
		Tuple:     results,
		Reprocess: reprocess,
	})
	if jsonResult.IsFailure() {
		log.Printf("[MPU][KPI] Serialization error: %s", jsonResult.GetError())
		return
	}
	err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.KPIFulfillmentCheckResultsQueueName), jsonResult.GetPayload())
	if err != nil {
		log.Printf("[MPU][KPI] Failed publishing KPI results: %s", err)
	}
}

func storeRaw(client rabbitmq.Client, results []sharedModel.TimeSeriesRawRecord) {
	if len(results) == 0 {
		return
	}
	jsonResult := sharedUtils.SerializeToJSON(results)
	if jsonResult.IsFailure() {
		log.Printf("[MPU][RAW] Serialization error: %s", jsonResult.GetError())
		return
	}
	err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesRawDataQueueName), jsonResult.GetPayload())
	if err != nil {
		log.Printf("[MPU][RAW] Failed storing RAW batch: %s", err)
	}
}

func storeKPI(client rabbitmq.Client, results []sharedModel.TimeSeriesKPIResultRecord) {
	if len(results) == 0 {
		return
	}
	jsonResult := sharedUtils.SerializeToJSON(results)
	if jsonResult.IsFailure() {
		log.Printf("[MPU][KPI] Serialization error: %s", jsonResult.GetError())
		return
	}
	err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesKPIResultQueueName), jsonResult.GetPayload())
	if err != nil {
		log.Printf("[MPU][KPI] Failed storing KPI batch: %s", err)
	}
}

func splitParamsBySDType(sdType string, params map[string]interface{}) (map[string]interface{}, map[string]string) {
	sdTypeDefinitionsMutex.RLock()
	definition := sdTypeDefinitions[sdType]
	sdTypeDefinitionsMutex.RUnlock()
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	for k, v := range params {
		paramDef, ok := definition.Parameters[k]
		if !ok {
			continue
		}
		if paramDef.Role == sharedModel.SDParameterRoleTag {
			tags[k] = fmt.Sprint(v)
			continue
		}
		fields[k] = v
	}
	return fields, tags
}

func mergeTagsAndFields(point sharedModel.TimeSeriesDataPoint) map[string]interface{} {
	params := make(map[string]interface{})
	for k, v := range point.Data {
		params[k] = v
	}
	for k, v := range point.Tags {
		switch k {
		case "sdInstanceUID", "sdType", "kpiDefinitionID":
			continue
		default:
			params[k] = v
		}
	}
	return params
}

func ProcessDelete(req sharedModel.KPIDeleteResultsRequestISCMessage) error {
	if req.KPIDefinitionID != 0 {
		key := makeKey(req.SDTypeUID, req.KPIDefinitionID)
		activeReprocessJobs.Delete(key)
		log.Printf("[MPU][REPROCESS] Delete request received -> cancelling active reprocess | kpiID=%d sdType=%s key=%s", req.KPIDefinitionID, req.SDTypeUID, key)
	}
	return nil
}

func makeKey(sdType string, kpiID uint32) string {
	return fmt.Sprintf("%s|%d", sdType, kpiID)
}

func isActive(jobKey string, jobID string) bool {
	v, ok := activeReprocessJobs.Load(jobKey)
	if !ok {
		return false
	}
	job, ok := v.(string)
	return ok && job == jobID
}

func waitForConfig(jobID string) {
	if jobID == "" {
		return
	}
	chIface, _ := jobSyncMap.LoadOrStore(jobID, make(chan struct{}))
	ch := chIface.(chan struct{})
	log.Printf("[MPU][SYNC] Waiting for config | jobID=%s", jobID)
	<-ch
	log.Printf("[MPU][SYNC] Config ready | jobID=%s", jobID)
}

func signalConfigReady(jobID string) {
	if jobID == "" {
		return
	}
	chIface, _ := jobSyncMap.LoadOrStore(jobID, make(chan struct{}))
	ch := chIface.(chan struct{})
	select {
	case <-ch:
	default:
		close(ch)
	}
	log.Printf("[MPU][SYNC] Config signaled | jobID=%s", jobID)
}
