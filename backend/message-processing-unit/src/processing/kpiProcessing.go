/**
 * @file kpiProcessing.go
 * @brief Zpracování KPI nad vstupními daty a příprava výsledků pro další vrstvy systému.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní vyhodnocování KPI a příprava výsledků k publikaci.
 * - Vojtěch Hubáček: doplnění stavové logiky a perzistence KPI výsledků do časové vrstvy.
 *
 * @ingroup riot_message_processing_unit
 */
package processing

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
)

type KPIProcessingOutputs struct {
	Results             []sharedModel.KPIFulfillmentCheckResultISCMessage
	TimeSeriesKPIRecord []sharedModel.TimeSeriesKPIResultRecord
}

func ProcessKPI(messagePayload sharedModel.KPIFulfillmentCheckRequestISCMessage, currentRaw map[string]interface{}, previousRaw map[string]interface{}, eventTime time.Time) KPIProcessingOutputs {
	kpiDefinitionsBySDTypeDenotationMapMutex.RLock()
	kpiDefinitions := kpiDefinitionsBySDTypeDenotationMap[messagePayload.SDTypeUID]
	kpiDefinitionsBySDTypeDenotationMapMutex.RUnlock()
	sdInstanceUID := messagePayload.SDInstanceUID
	results := []sharedModel.KPIFulfillmentCheckResultISCMessage{}
	tsRecords := []sharedModel.TimeSeriesKPIResultRecord{}
	_, tags := splitParamsBySDType(messagePayload.SDTypeUID, currentRaw)
	context := KPIEvaluationContext{
		CurrentValues:  currentRaw,
		PreviousValues: previousRaw,
	}
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
		result := CheckKPIFulfillment(kpiDefinition, context)
		if result.IsFailure() {
			continue
		}
		value := result.GetPayload()
		kpiUID := kpiDefinition.UID
		key := sharedModel.KPIKey{
			SDInstanceUID:    sdInstanceUID,
			KPIDefinitionUID: kpiUID,
		}
		last, exists := lastKPI.Load(key)
		if exists {
			state := last.(sharedModel.KPIState)
			if !eventTime.Before(state.EventTime) && state.Value == value {
				state.SynchronizedAt = time.Now()
				lastKPI.Store(key, state)
				continue
			}
		}
		lastKPI.Store(key, sharedModel.KPIState{
			Value:          value,
			EventTime:      eventTime,
			SynchronizedAt: time.Now(),
		})
		results = append(results, sharedModel.KPIFulfillmentCheckResultISCMessage{
			SDTypeUID:        messagePayload.SDTypeUID,
			EventTime:        eventTime,
			SDInstanceUID:    sdInstanceUID,
			KPIDefinitionUID: kpiUID,
			Fulfilled:        value,
		})
		toStore := sharedModel.TimeSeriesKPIResultRecord{
			JobID:            "",
			EventTime:        eventTime,
			SDInstanceUID:    sdInstanceUID,
			SDTypeUID:        messagePayload.SDTypeUID,
			KPIDefinitionUID: kpiUID,
			Fulfilled:        value,
			Tags:             tags,
		}
		tsRecords = append(tsRecords, toStore)
	}
	return KPIProcessingOutputs{
		Results:             results,
		TimeSeriesKPIRecord: tsRecords,
	}
}
