/**
 * @file kpiFulfillmentCheckResults.go
 * @brief Mapování výsledků vyhodnocení KPI z doménového modelu do GraphQL modelu.
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
package dll2gql

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelKPIFulfillmentCheckResult(kpiFulfillmentCheckResult dllModel.KPIFulfillmentCheckResult) graphQLModel.KPIFulfillmentCheckResult {
	return graphQLModel.KPIFulfillmentCheckResult{
		SdTypeUID:        kpiFulfillmentCheckResult.SDTypeUID,
		SdInstanceUID:    kpiFulfillmentCheckResult.SDInstanceUID,
		KpiDefinitionUID: kpiFulfillmentCheckResult.KPIDefinitionUID,
		Fulfilled:        kpiFulfillmentCheckResult.Fulfilled,
		EventTime:        kpiFulfillmentCheckResult.EventTime.Format(time.RFC3339Nano),
	}
}
