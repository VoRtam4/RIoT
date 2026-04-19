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
