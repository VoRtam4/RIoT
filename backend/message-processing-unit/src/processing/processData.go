package processing

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	kpiDefinitionsBySDTypeDenotationMap      map[string][]sharedModel.KPIDefinition = make(map[string][]sharedModel.KPIDefinition)
	kpiDefinitionsBySDTypeDenotationMapMutex sync.RWMutex
	unitUUID                                 string
	limit                                    = 500
	sdTypeDefinitions                        = map[string]map[string]sharedModel.SDParameter{}
	sdTypeDefinitionsMutex                   sync.RWMutex
)

func InitializeProcessing() {
	unitUUID = uuid.New().String()
	log.Printf("[MPU] Processing unit initialized | unitUUID=%s", unitUUID)
}

func DenotationMapUpdates(rabbitMQClient rabbitmq.Client) {
	log.Printf("[MPU] Waiting for KPI configuration updates | unit=%s", unitUUID)
	err := rabbitmq.ConsumeJSONMessagesFromFanoutExchange[sharedModel.KPIConfigurationUpdateISCMessage](rabbitMQClient, sharedConstants.BuiltInFanoutExchangeName, func(messagePayload sharedModel.KPIConfigurationUpdateISCMessage) error {
		log.Printf("[MPU] KPI configuration update received | unit=%s entries=%d", unitUUID, len(messagePayload))
		kpiDefinitionsBySDTypeDenotationMapMutex.Lock()
		kpiDefinitionsBySDTypeDenotationMap = messagePayload
		kpiDefinitionsBySDTypeDenotationMapMutex.Unlock()
		log.Printf("[MPU] KPI configuration map updated successfully | unit=%s", unitUUID)
		return nil
	})
	if err != nil {
		log.Printf("[MPU] Consumption of messages from the '%s' fanout exchange has failed: %s\n", sharedConstants.BuiltInFanoutExchangeName, err.Error())
	}
}

func UpdateSDType(messages []sharedModel.SDTypeRegistrationRequestISCMessage) {
	log.Printf("[MPU][SDTYPE] Received SDType update | count=%d", len(messages))
	newMap := make(map[string]map[string]sharedModel.SDParameter)
	for _, msg := range messages {
		paramMap := make(map[string]sharedModel.SDParameter)
		for _, p := range msg.Parameters {
			paramMap[p.Denotation] = p
		}
		newMap[msg.SDTypeSpecification] = paramMap
	}
	sdTypeDefinitionsMutex.Lock()
	sdTypeDefinitions = newMap
	sdTypeDefinitionsMutex.Unlock()
	log.Printf("[MPU][SDTYPE] SDType cache updated | types=%d", len(newMap))
}

func ProcessRaw(rabbitMQClient rabbitmq.Client, messagePayload sharedModel.KPIFulfillmentCheckRequestISCMessage, params map[string]interface{}, eventTime time.Time) bool {
	log.Printf("[MPU][RAW] Processing raw data | uid=%s type=%s time=%v params=%v", messagePayload.SDInstanceUID, messagePayload.SDTypeSpecification, eventTime, params)
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
	changed := false
	for k, v := range params {
		oldVal, ok := lastSnapshot.Values[k]
		if !ok {
			changed = true
			break
		}
		if !sharedUtils.CompareJSONs(oldVal, v) {
			changed = true
			break
		}
	}
	if !changed {
		for k := range lastSnapshot.Values {
			if _, ok := params[k]; !ok {
				changed = true
				break
			}
		}
	}
	if exists {
		lastSnapshot.SynchronizedAt = time.Now()
		lastRaw.Store(uid, lastSnapshot)
	}
	newValues := make(map[string]interface{}, len(params))
	for k, v := range params {
		newValues[k] = v
	}
	for k := range lastSnapshot.Values {
		if _, ok := params[k]; !ok {
			newValues[k] = nil
		}
	}
	fields, tags := splitParamsBySDType(messagePayload.SDTypeSpecification, newValues)
	if exists && eventTime.Before(lastSnapshot.EventTime) {
		log.Printf("[MPU][RAW] Older event -> storing ONLY TS | uid=%s", uid)
		storeRaw(rabbitMQClient, []sharedModel.TimeSeriesRawRecord{{
			EventTime:           eventTime,
			SDInstanceUID:       uid,
			SDTypeSpecification: messagePayload.SDTypeSpecification,
			Fields:              fields,
			Tags:                tags,
		}})
		return true
	}
	if !changed {
		log.Printf("[MPU][RAW] No change detected -> skipping RAW write | uid=%s", uid)
		return false
	}
	newSnapshot := sharedModel.RawState{
		Values:         newValues,
		EventTime:      eventTime,
		SynchronizedAt: time.Now(),
	}
	lastRaw.Store(uid, newSnapshot)
	log.Printf("[MPU][RAW] Change detected -> storing RAW snapshot | uid=%s snapshot=%v", uid, newValues)
	storeRaw(rabbitMQClient, []sharedModel.TimeSeriesRawRecord{{
		EventTime:           eventTime,
		SDInstanceUID:       uid,
		SDTypeSpecification: messagePayload.SDTypeSpecification,
		Fields:              fields,
		Tags:                tags,
	}})
	return true
}

func ProcessKPI(rabbitMQClient rabbitmq.Client, messagePayload sharedModel.KPIFulfillmentCheckRequestISCMessage, params map[string]interface{}, eventTime time.Time) error {
	log.Printf("[MPU][KPI] Processing KPI evaluation | uid=%s type=%s time=%v", messagePayload.SDInstanceUID, messagePayload.SDTypeSpecification, eventTime)
	kpiDefinitionsBySDTypeDenotationMapMutex.RLock()
	kpiDefinitions := kpiDefinitionsBySDTypeDenotationMap[messagePayload.SDTypeSpecification]
	kpiDefinitionsBySDTypeDenotationMapMutex.RUnlock()
	log.Printf("[MPU][KPI] Found KPI definitions | count=%d type=%s", len(kpiDefinitions), messagePayload.SDTypeSpecification)
	sdInstanceUID := messagePayload.SDInstanceUID
	results := sharedUtils.EmptySlice[sharedModel.KPIFulfillmentCheckResultISCMessage]()
	tsRecords := sharedUtils.EmptySlice[sharedModel.TimeSeriesKPIResultRecord]()
	_, tags := splitParamsBySDType(messagePayload.SDTypeSpecification, params)
	var paramAny any = params
	for _, kpiDefinition := range kpiDefinitions {
		log.Printf("[MPU][KPI] Evaluating KPI | id=%v uid=%s", kpiDefinition.ID, sdInstanceUID)
		if kpiDefinition.SDInstanceMode != sharedModel.ALL {
			containsUID := false
			for _, uid := range kpiDefinition.SelectedSDInstanceUIDs {
				if uid == sdInstanceUID {
					containsUID = true
					break
				}
			}
			if !containsUID {
				log.Printf("[MPU][KPI] Instance not selected for KPI | kpiID=%v uid=%s", kpiDefinition.ID, sdInstanceUID)
				continue
			}
		}
		result := CheckKPIFulfillment(kpiDefinition, &paramAny)
		if result.IsFailure() {
			log.Printf("[MPU][KPI] KPI evaluation failed | kpiID=%v error=%s", kpiDefinition.ID, result.GetError().Error())
			continue
		}
		value := result.GetPayload()
		kpiID := sharedUtils.NewOptionalFromPointer(kpiDefinition.ID).GetPayload()
		log.Printf("[MPU][KPI] KPI evaluation result | kpiID=%d uid=%s fulfilled=%v", kpiID, sdInstanceUID, value)
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
			EventTime:       eventTime,
			SDInstanceUID:   sdInstanceUID,
			KPIDefinitionID: kpiID,
			Fulfilled:       value,
		})
		tsRecords = append(tsRecords, sharedModel.TimeSeriesKPIResultRecord{
			EventTime:           eventTime,
			SDInstanceUID:       sdInstanceUID,
			SDTypeSpecification: messagePayload.SDTypeSpecification,
			KPIDefinitionID:     kpiID,
			Fulfilled:           value,
			Tags:                tags,
		})
	}
	log.Printf("[MPU][KPI] Storing KPI results to TS | count=%d", len(tsRecords))
	storeKPI(rabbitMQClient, tsRecords)
	log.Printf("[MPU][KPI] Publishing KPI results | count=%d", len(results))
	publishKPI(rabbitMQClient, results)
	return nil
}

func ReprocessKPI(req sharedModel.KPIReprocessRequestISCMessage) error {
	log.Printf("[MPU][REPROCESS] Starting KPI reprocess | kpiID=%d type=%s until=%v", req.KPIDefinitionID, req.SDTypeSpecification, req.To)
	client := rabbitmq.NewClient()
	defer client.Dispose()
	kpiDefinitionsBySDTypeDenotationMapMutex.RLock()
	kpis := kpiDefinitionsBySDTypeDenotationMap[req.SDTypeSpecification]
	kpiDefinitionsBySDTypeDenotationMapMutex.RUnlock()
	var targetKPI *sharedModel.KPIDefinition
	for i := range kpis {
		if kpis[i].ID != nil && *kpis[i].ID == req.KPIDefinitionID {
			targetKPI = &kpis[i]
			break
		}
	}
	if targetKPI == nil {
		return fmt.Errorf("KPI definition not found")
	}
	log.Printf("[MPU][REPROCESS] KPI definition located | id=%d", req.KPIDefinitionID)
	readReq := sharedModel.TimeSeriesReprocessReadRequest{
		SDTypeSpecification: req.SDTypeSpecification,
		SDInstanceUIDs:      req.SDInstanceUIDs,
		To:                  req.To,
		Batch:               limit,
	}
	jsonReq := sharedUtils.SerializeToJSON(readReq)
	if jsonReq.IsFailure() {
		return jsonReq.GetError()
	}
	correlationID := uuid.New().String()
	err := client.PublishJSONMessageRPC(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesReprocessReadRequestQueueName), jsonReq.GetPayload(), correlationID, sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesReprocessReadResponseQueueName))
	if err != nil {
		return err
	}
	instanceState := map[string]sharedModel.KPIState{}
	handler := createReprocessHandler(client, req, targetKPI, instanceState)
	return rabbitmq.ConsumeJSONMessagesWithAccessToDelivery[sharedModel.TimeSeriesReprocessReadResponse](client, sharedConstants.TimeSeriesReprocessReadResponseQueueName, correlationID, handler)
}

func createReprocessHandler(client rabbitmq.Client, req sharedModel.KPIReprocessRequestISCMessage, targetKPI *sharedModel.KPIDefinition, instanceState map[string]sharedModel.KPIState) func(resp sharedModel.TimeSeriesReprocessReadResponse, delivery amqp.Delivery) error {
	return func(resp sharedModel.TimeSeriesReprocessReadResponse, delivery amqp.Delivery) error {
		log.Printf("[MPU][REPROCESS] Batch received | points=%d hasMore=%v", len(resp.Data), resp.HasMore)
		if resp.Error != "" {
			return fmt.Errorf(resp.Error)
		}
		tsBatch := make([]sharedModel.TimeSeriesKPIResultRecord, 0, len(resp.Data))
		for _, point := range resp.Data {
			result, ok := evaluateReprocessPoint(point, req.KPIDefinitionID, targetKPI)
			if !ok {
				continue
			}
			prev, exists := instanceState[result.SDInstanceUID]
			if exists && prev.Value == result.Fulfilled {
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
				EventTime:           result.EventTime,
				SDInstanceUID:       result.SDInstanceUID,
				SDTypeSpecification: req.SDTypeSpecification,
				KPIDefinitionID:     result.KPIDefinitionID,
				Fulfilled:           result.Fulfilled,
				Tags:                tags,
			})
		}
		if len(tsBatch) > 0 {
			log.Printf("[MPU][REPROCESS] Writing KPI batch | size=%d", len(tsBatch))
			storeKPI(client, tsBatch)
		}
		if !resp.HasMore {
			log.Printf("[MPU][REPROCESS] Finalizing reprocess")
			finalizeReprocess(client, req, instanceState)
			return nil
		}
		return nil
	}
}

func evaluateReprocessPoint(point sharedModel.TimeSeriesDataPoint, kpiID uint32, targetKPI *sharedModel.KPIDefinition) (sharedModel.KPIFulfillmentCheckResultISCMessage, bool) {
	log.Printf("[MPU][REPROCESS] Evaluating RAW point for KPI | time=%v tags=%v", point.Time, point.Tags)
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
	log.Printf("[MPU][REPROCESS] KPI evaluated | uid=%s fulfilled=%v", sdInstanceUID, result.GetPayload())
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
		last, exists := lastKPI.Load(key)
		if exists {
			lastState := last.(sharedModel.KPIState)
			if !lastState.EventTime.Before(req.To) {
				continue
			}
		}
		lastKPI.Store(key, state)
		results = append(results, sharedModel.KPIFulfillmentCheckResultISCMessage{
			EventTime:       state.EventTime,
			SDInstanceUID:   uid,
			KPIDefinitionID: req.KPIDefinitionID,
			Fulfilled:       state.Value,
		})
		if len(results) >= limit {
			publishKPI(client, results)
			results = results[:0]
		}
	}
	if len(results) > 0 {
		publishKPI(client, results)
	}
}

func publishKPI(client rabbitmq.Client, results []sharedModel.KPIFulfillmentCheckResultISCMessage) {
	log.Printf("[MPU][KPI] publishKPI called | results=%d", len(results))
	if len(results) == 0 {
		return
	}
	jsonResult := sharedUtils.SerializeToJSON(results)
	if jsonResult.IsFailure() {
		log.Printf("[MPU][KPI] Serialization error: %s", jsonResult.GetError())
		return
	}
	err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.KPIFulfillmentCheckResultsQueueName), jsonResult.GetPayload())
	if err != nil {
		log.Printf("[MPU][KPI] Failed publishing KPI results: %s", err)
	}
	log.Printf("[MPU][KPI] KPI results published successfully")
}

func storeRaw(client rabbitmq.Client, results []sharedModel.TimeSeriesRawRecord) {
	log.Printf("[MPU][RAW] storeRaw called | count=%d", len(results))
	for _, result := range results {
		jsonResult := sharedUtils.SerializeToJSON(result)
		if jsonResult.IsFailure() {
			log.Printf("[MPU][RAW] Serialization error: %s", jsonResult.GetError())
			continue
		}
		err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesRawDataQueueName), jsonResult.GetPayload())
		if err != nil {
			log.Printf("[MPU][RAW] Failed storing RAW record: %s", err)
		}
	}
}

func storeKPI(client rabbitmq.Client, results []sharedModel.TimeSeriesKPIResultRecord) {
	log.Printf("[MPU][KPI] storeKPI called | count=%d", len(results))
	for _, result := range results {
		jsonResult := sharedUtils.SerializeToJSON(result)
		if jsonResult.IsFailure() {
			log.Printf("[MPU][KPI] Serialization error: %s", jsonResult.GetError())
			continue
		}
		err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesKPIResultQueueName), jsonResult.GetPayload())
		if err != nil {
			log.Printf("[MPU][KPI] Failed storing KPI record: %s", err)
		}
	}
}

func splitParamsBySDType(sdType string, params map[string]interface{}) (map[string]interface{}, map[string]string) {
	sdTypeDefinitionsMutex.RLock()
	definition := sdTypeDefinitions[sdType]
	sdTypeDefinitionsMutex.RUnlock()
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	for k, v := range params {
		role := definition[k].Role
		if role == "TAG" {
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
