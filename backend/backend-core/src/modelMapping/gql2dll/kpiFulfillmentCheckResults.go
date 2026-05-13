/**
 * @file kpiFulfillmentCheckResults.go
 * @brief Mapování GraphQL požadavků na výsledky KPI do doménového modelu.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package gql2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToDLLModelKPIFulfillmentCheckResultRequest(kpiFulfillmentCheckResultRequest graphQLModel.KPIFulfillmentCheckResultRequest) dllModel.KPIFulfillmentCheckResultRequest {
	return dllModel.KPIFulfillmentCheckResultRequest{
		KPIDefinitionID: kpiFulfillmentCheckResultRequest.KpiDefinitionID,
		SDInstanceID:    kpiFulfillmentCheckResultRequest.SdInstanceID,
	}
}
