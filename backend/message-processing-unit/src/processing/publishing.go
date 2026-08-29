/**
 * @file publishing.go
 * @brief Publikování výsledků zpracování, aktualizací typů zdrojů a dat pro časovou vrstvu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní publikování surových dat a KPI výsledků.
 * - Vojtěch Hubáček: rozšíření publikace o batching a implementace ostatních publikačních toků souboru.
 *
 * @ingroup riot_message_processing_unit
 */
package processing

import (
	"fmt"
	"log"
	"strings"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/rabbitmq"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedConstants"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func publishSDTypeUpdate(client rabbitmq.Client, sdTypeUID string) {
	publishSDTypeUpdates(client, []string{sdTypeUID})
}

func publishSDTypeUpdates(client rabbitmq.Client, sdTypeUIDs []string) {
	if len(sdTypeUIDs) == 0 {
		return
	}
	messages := make([]sharedModel.SDTypeRegistrationRequestISCMessage, 0, len(sdTypeUIDs))
	seen := make(map[string]struct{}, len(sdTypeUIDs))
	for _, sdTypeUID := range sdTypeUIDs {
		if _, exists := seen[sdTypeUID]; exists {
			continue
		}
		seen[sdTypeUID] = struct{}{}
		sdTypeDefinitionsMutex.RLock()
		definition, exists := sdTypeDefinitions[sdTypeUID]
		sdTypeDefinitionsMutex.RUnlock()
		if !exists {
			log.Printf("[MPU][SDTYPE] Definition not found for type=%s", sdTypeUID)
			continue
		}
		params := make([]sharedModel.SDParameter, 0, len(definition.Parameters))
		for _, p := range definition.Parameters {
			params = append(params, p)
		}
		label := definition.Label
		if label == "" {
			label = strings.ToUpper(sdTypeUID[:1]) + sdTypeUID[1:]
		}
		messages = append(messages, sharedModel.SDTypeRegistrationRequestISCMessage{
			SDTypeUID:  sdTypeUID,
			Label:      label,
			Parameters: params,
		})
	}
	if len(messages) == 0 {
		return
	}
	err := rabbitmq.PublishJSONBatches(
		client,
		sharedUtils.NewEmptyOptional[string](),
		sharedUtils.NewOptionalOf(sharedConstants.SDTypeRegistrationRequestsQueueName),
		messages,
		limit,
	)
	if err != nil {
		log.Printf("[MPU][SDTYPE] Failed publishing SDType registration/upsert tuple: %s", err)
	}
}

func publishRaw(client rabbitmq.Client, results []sharedModel.RawDataPointISCMessage) {
	if len(results) == 0 {
		return
	}
	err := rabbitmq.PublishJSONBatches(
		client,
		sharedUtils.NewEmptyOptional[string](),
		sharedUtils.NewOptionalOf(sharedConstants.RawDataPointQueueName),
		results,
		limit,
	)
	if err != nil {
		log.Printf("[MPU][RAW] Failed publishing RAW data: %s", err)
		return
	}
}

func publishKPI(client rabbitmq.Client, results []sharedModel.KPIFulfillmentCheckResultISCMessage, reprocess bool) {
	if len(results) == 0 {
		return
	}
	for start := 0; start < len(results); start += limit {
		end := start + limit
		if end > len(results) {
			end = len(results)
		}
		jsonResult := sharedUtils.SerializeToJSON(sharedModel.KPIFulfillmentCheckResultTupleISCMessage{
			Tuple:     results[start:end],
			Reprocess: reprocess,
		})
		if jsonResult.IsFailure() {
			log.Printf("[MPU][KPI] Serialization error: %s", jsonResult.GetError())
			return
		}
		err := client.PublishJSONMessage(sharedUtils.NewEmptyOptional[string](), sharedUtils.NewOptionalOf(sharedConstants.KPIFulfillmentCheckResultsQueueName), jsonResult.GetPayload())
		if err != nil {
			log.Printf("[MPU][KPI] Failed publishing KPI results: %s", err)
			return
		}
	}
}

func storeRaw(client rabbitmq.Client, results []sharedModel.TimeSeriesRawRecord) {
	if len(results) == 0 {
		return
	}
	err := rabbitmq.PublishJSONBatches(
		client,
		sharedUtils.NewEmptyOptional[string](),
		sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesRawDataQueueName),
		results,
		limit,
	)
	if err != nil {
		log.Printf("[MPU][RAW] Failed storing RAW batch: %s", err)
	}
}

func storeKPI(client rabbitmq.Client, results []sharedModel.TimeSeriesKPIResultRecord) {
	if len(results) == 0 {
		return
	}
	err := rabbitmq.PublishJSONBatches(
		client,
		sharedUtils.NewEmptyOptional[string](),
		sharedUtils.NewOptionalOf(sharedConstants.TimeSeriesKPIResultQueueName),
		results,
		limit,
	)
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
		case "sdInstanceUID", "sdType", "kpiDefinitionUID":
			continue
		default:
			params[k] = v
		}
	}
	return params
}
