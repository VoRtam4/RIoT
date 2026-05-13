/**
 * @file kpiFulfillmentCheckResults.go
 * @brief Mapování výsledků vyhodnocení KPI z doménového modelu do databázového modelu.
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
package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDBModelEntityKPIFulfillmentCheckResult(kpi dllModel.KPIFulfillmentCheckResult) dbModel.KPIFulfillmentCheckResultEntity {
	return dbModel.KPIFulfillmentCheckResultEntity{
		SDInstanceID:    kpi.SDInstanceID,
		KPIDefinitionID: kpi.KPIDefinitionID,
		Fulfilled:       kpi.Fulfilled,
		EventTime:       kpi.EventTime,
	}
}
