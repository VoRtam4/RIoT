/**
 * @file rawProcessing.go
 * @brief Normalizace příchozích surových dat a příprava změn pro další části platformy.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní základ zpracování příchozích surových hodnot.
 * - Vojtěch Hubáček: doplnění práce se stavem, časovou vrstvou, kontrol parametrů a aktualizací typů zdrojů.
 *
 * @ingroup riot_message_processing_unit
 */
package processing

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type RawProcessingOutputs struct {
	SDTypeChanged bool
	RawMessages   []sharedModel.RawDataPointISCMessage
	RawRecords    []sharedModel.TimeSeriesRawRecord
}

func ProcessRaw(messagePayload sharedModel.KPIFulfillmentCheckRequestISCMessage, params map[string]interface{}, eventTime time.Time) (RawProcessingOutputs, bool) {
	uid := messagePayload.SDInstanceUID
	outputs := RawProcessingOutputs{}
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
	outputs.SDTypeChanged = result.SDTypeChanged
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
		outputs.RawRecords = append(outputs.RawRecords, sharedModel.TimeSeriesRawRecord{
			EventTime:     eventTime,
			SDInstanceUID: uid,
			SDTypeUID:     messagePayload.SDTypeUID,
			Fields:        fields,
			Tags:          tags,
		})
		return outputs, true
	}
	if !changed {
		return outputs, false
	}
	payloadResult := sharedUtils.SerializeToJSON(newValues)
	if payloadResult.IsFailure() {
		log.Printf("[MPU][RAW] Payload serialization error: %s", payloadResult.GetError())
		return outputs, false
	}
	outputs.RawMessages = append(outputs.RawMessages, sharedModel.RawDataPointISCMessage{
		SDTypeUID:     messagePayload.SDTypeUID,
		SDInstanceUID: uid,
		EventTime:     eventTime,
		Payload:       payloadResult.GetPayload(),
	})
	lastRaw.Store(uid, sharedModel.RawState{
		Values:         newValues,
		EventTime:      eventTime,
		SynchronizedAt: time.Now(),
	})
	outputs.RawRecords = append(outputs.RawRecords, sharedModel.TimeSeriesRawRecord{
		EventTime:     eventTime,
		SDInstanceUID: uid,
		SDTypeUID:     messagePayload.SDTypeUID,
		Fields:        fields,
		Tags:          tags,
	})
	return outputs, true
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
