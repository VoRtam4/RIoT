/**
 * @file kpiFulfillmentCheckResults.go
 * @brief Mapování výsledků vyhodnocení KPI z databázového modelu do doménového modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování výsledků vyhodnocení KPI.
 * - Vojtěch Hubáček: doplnění mapování času události EventTime.
 *
 * @ingroup riot_backend_core
 */
package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDLLModelKPIFulfillmentCheckResult(row dbModel.KPIFulfillmentCheckResultEntity) dllModel.KPIFulfillmentCheckResult {
	return dllModel.KPIFulfillmentCheckResult{
		SDTypeID:         row.SDInstance.SDTypeID,
		SDTypeUID:        row.SDInstance.SDType.UID,
		SDInstanceID:     row.SDInstanceID,
		SDInstanceUID:    row.SDInstance.UID,
		KPIDefinitionID:  row.KPIDefinitionID,
		KPIDefinitionUID: kpiDefinitionUID(row.KPIDefinition),
		Fulfilled:        row.Fulfilled,
		EventTime:        row.EventTime,
	}
}

func kpiDefinitionUID(kpiDefinition dbModel.KPIDefinitionEntity) string {
	if kpiDefinition.UID == nil {
		return ""
	}
	return *kpiDefinition.UID
}
