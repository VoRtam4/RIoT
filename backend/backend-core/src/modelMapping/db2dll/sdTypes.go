/**
 * @file sdTypes.go
 * @brief Mapování typů zdrojů dat z databázového modelu do doménového modelu.
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
