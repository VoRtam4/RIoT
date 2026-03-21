package dll2gql

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToGraphQLModelSDType(sdType dllModel.SDType) graphQLModel.SDType {
	return graphQLModel.SDType{
		ID:         sdType.ID.GetPayload(),
		Denotation: sdType.Denotation,
		Parameters: sharedUtils.Map(sdType.Parameters, func(sdParameter dllModel.SDParameter) graphQLModel.SDParameter {
			return graphQLModel.SDParameter{
				ID:         sdParameter.ID.GetPayload(),
				Denotation: sdParameter.Denotation,
				Type: func(sdParameterType dllModel.SDParameterType) graphQLModel.SDParameterType {
					switch sdParameterType {
					case dllModel.SDParameterTypeString:
						return graphQLModel.SDParameterTypeString
					case dllModel.SDParameterTypeNumber:
						return graphQLModel.SDParameterTypeNumber
					case dllModel.SDParameterTypeBoolean:
						return graphQLModel.SDParameterTypeBoolean
					}
					panic(fmt.Errorf("unpexted model mapping failure – shouldn't happen"))
				}(sdParameter.Type),
				Role: func(sdParameterRole dllModel.SDParameterRole) graphQLModel.SDParameterRole {
					switch sdParameterRole {
					case dllModel.SDParameterRoleTag:
						return graphQLModel.SDParameterRoleTag
					case dllModel.SDParameterRoleField:
						return graphQLModel.SDParameterRoleField
					}
					panic(fmt.Errorf("unpexted model mapping failure – shouldn't happen"))
				}(sdParameter.Role),
			}
		}),
	}
}
