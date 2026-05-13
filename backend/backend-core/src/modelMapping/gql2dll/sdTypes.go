/**
 * @file sdTypes.go
 * @brief Mapování vstupů typů zdrojů dat z GraphQL modelu do doménového modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování typů zdrojů dat a jejich parametrů.
 * - Vojtěch Hubáček: doplnění mapování labelů, rolí parametrů field/tag a zjednodušení převodu string enumů přetypováním na doménové typy.
 *
 * @ingroup riot_backend_core
 */
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
