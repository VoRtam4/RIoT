package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDLLModelRawDataPoint(rawDataPointEntity dbModel.RawDataPointEntity) dllModel.RawDataPoint {
	return dllModel.RawDataPoint{
		SDTypeID:     rawDataPointEntity.SDInstance.SDTypeID,
		SDInstanceID: rawDataPointEntity.SDInstanceID,
		EventTime:    rawDataPointEntity.EventTime,
		Payload:      rawDataPointEntity.Payload,
	}
}
