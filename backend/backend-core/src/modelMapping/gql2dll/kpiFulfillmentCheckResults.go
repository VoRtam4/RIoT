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
