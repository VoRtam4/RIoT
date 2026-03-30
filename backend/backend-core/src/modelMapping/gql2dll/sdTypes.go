package gql2dll

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelSDType(sdTypeInput graphQLModel.SDTypeInput) dllModel.SDType {
	return dllModel.SDType{
		ID:         sharedUtils.NewEmptyOptional[uint32](),
		Label:      sdTypeInput.Label,
		Denotation: sdTypeInput.Denotation,
		Parameters: sharedUtils.Map(sdTypeInput.Parameters, func(sdParameterInput graphQLModel.SDParameterInput) dllModel.SDParameter {
			return dllModel.SDParameter{
				ID:         sharedUtils.NewEmptyOptional[uint32](),
				Label:      sdParameterInput.Label,
				Denotation: sdParameterInput.Denotation,
				Type: func(sdParameterType graphQLModel.SDParameterType) dllModel.SDParameterType {
					switch sdParameterType {
					case graphQLModel.SDParameterTypeString:
						return dllModel.SDParameterTypeString
					case graphQLModel.SDParameterTypeNumber:
						return dllModel.SDParameterTypeNumber
					case graphQLModel.SDParameterTypeBoolean:
						return dllModel.SDParameterTypeBoolean
					}
					panic(fmt.Errorf("unpexted model mapping failure – shouldn't happen"))
				}(sdParameterInput.Type),
				Role: func(sdParameterRole graphQLModel.SDParameterRole) dllModel.SDParameterRole {
					switch sdParameterRole {
					case graphQLModel.SDParameterRoleTag:
						return dllModel.SDParameterRoleTag
					case graphQLModel.SDParameterRoleField:
						return dllModel.SDParameterRoleField
					}
					panic(fmt.Errorf("unpexted model mapping failure – shouldn't happen"))
				}(sdParameterInput.Role),
			}
		}),
	}
}
