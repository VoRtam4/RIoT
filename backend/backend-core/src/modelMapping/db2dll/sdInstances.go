/**
 * @file sdInstances.go
 * @brief Mapování instancí zdrojů dat z databázového modelu do doménového modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování instancí zdrojů dat.
 * - Vojtěch Hubáček: doplnění mapování labelů.
 *
 * @ingroup riot_backend_core
 */
package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelSDInstance(sdInstanceEntity dbModel.SDInstanceEntity) dllModel.SDInstance {
	return dllModel.SDInstance{
		ID:              sharedUtils.NewOptionalOf[uint32](sdInstanceEntity.ID),
		UID:             sdInstanceEntity.UID,
		Label:           sdInstanceEntity.Label,
		ConfirmedByUser: sdInstanceEntity.ConfirmedByUser,
		UserIdentifier:  sdInstanceEntity.UserIdentifier,
		SDType:          ToDLLModelSDType(sdInstanceEntity.SDType),
	}
}
