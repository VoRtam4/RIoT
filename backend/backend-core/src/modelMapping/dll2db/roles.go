/**
 * @file roles.go
 * @brief Mapování rolí z doménového modelu do databázového modelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDBModelRole(role dllModel.Role) dbModel.RoleEntity {
	return dbModel.RoleEntity{
		ID:     role.ID,
		UID:    role.UID,
		Label:  role.Label,
		System: role.System,
	}
}
