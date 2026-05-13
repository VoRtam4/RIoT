/**
 * @file sdInstances.go
 * @brief Mapování instancí zdrojů dat z doménového modelu do databázového modelu.
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
package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDBModelEntitySDInstance(sdInstance dllModel.SDInstance) dbModel.SDInstanceEntity {
	return dbModel.SDInstanceEntity{
		ID:              sdInstance.ID.GetPayloadOrDefault(0),
		UID:             sdInstance.UID,
		Label:           sdInstance.Label,
		ConfirmedByUser: sdInstance.ConfirmedByUser,
		UserIdentifier:  sdInstance.UserIdentifier,
		SDTypeID:        sdInstance.SDType.ID.GetPayloadOrDefault(0),
	}
}
