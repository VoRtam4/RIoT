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
		SDTypeID:        row.SDInstance.SDTypeID,
		SDInstanceID:    row.SDInstanceID,
		KPIDefinitionID: row.KPIDefinitionID,
		Fulfilled:       row.Fulfilled,
		EventTime:       row.EventTime,
	}
}
