/**
 * @file sdTypes.go
 * @brief Mapování typů zdrojů dat z doménového modelu do GraphQL modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování typů zdrojů dat a jejich parametrů.
 * - Vojtěch Hubáček: doplnění mapování labelů, rolí parametrů field/tag a zjednodušení převodu string enumů přetypováním na GraphQL typy.
 *
 * @ingroup riot_backend_core
 */
package dll2gql

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToGraphQLModelSDType(sdType dllModel.SDType) graphQLModel.SDType {
	return graphQLModel.SDType{
		ID:    sdType.ID.GetPayload(),
		UID:   sdType.UID,
		Label: sdType.Label,
		Parameters: sharedUtils.Map(sdType.Parameters, func(sdParameter dllModel.SDParameter) graphQLModel.SDParameter {
			return graphQLModel.SDParameter{
				ID:         sdParameter.ID.GetPayload(),
				Label:      sdParameter.Label,
				Denotation: sdParameter.Denotation,
				Type:       graphQLModel.SDParameterType(sdParameter.Type),
				Role:       graphQLModel.SDParameterRole(sdParameter.Role),
			}
		}),
	}
}
