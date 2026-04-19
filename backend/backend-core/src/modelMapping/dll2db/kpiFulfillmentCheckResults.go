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
