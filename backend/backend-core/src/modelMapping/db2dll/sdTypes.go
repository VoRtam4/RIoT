package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelSDType(sdTypeEntity dbModel.SDTypeEntity) dllModel.SDType {
	return dllModel.SDType{
		ID:    sharedUtils.NewOptionalOf[uint32](sdTypeEntity.ID),
		UID:   sdTypeEntity.UID,
		Label: sdTypeEntity.Label,
		Parameters: sharedUtils.Map(sdTypeEntity.Parameters, func(sdParameterEntity dbModel.SDParameterEntity) dllModel.SDParameter {
			return dllModel.SDParameter{
				ID:         sharedUtils.NewOptionalOf[uint32](sdParameterEntity.ID),
				Label:      sdParameterEntity.Label,
				Denotation: sdParameterEntity.Denotation,
				Type:       dllModel.SDParameterType(sdParameterEntity.Type),
				Role:       dllModel.SDParameterRole(sdParameterEntity.Role),
			}
		}),
	}
}
