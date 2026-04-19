package gql2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelSDType(sdTypeInput graphQLModel.SDTypeInput) dllModel.SDType {
	return dllModel.SDType{
		ID:    sharedUtils.NewEmptyOptional[uint32](),
		Label: sdTypeInput.Label,
		UID:   sdTypeInput.UID,
		Parameters: sharedUtils.Map(sdTypeInput.Parameters, func(sdParameterInput graphQLModel.SDParameterInput) dllModel.SDParameter {
			return dllModel.SDParameter{
				ID:         sharedUtils.NewEmptyOptional[uint32](),
				Label:      sdParameterInput.Label,
				Denotation: sdParameterInput.Denotation,
				Type:       dllModel.SDParameterType(sdParameterInput.Type),
				Role:       dllModel.SDParameterRole(sdParameterInput.Role),
			}
		}),
	}
}
