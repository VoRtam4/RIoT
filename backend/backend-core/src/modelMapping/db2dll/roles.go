/**
 * @file roles.go
 * @brief Mapování rolí a oprávnění z databázového modelu do doménového modelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelRole(role dbModel.RoleEntity) dllModel.Role {
	return dllModel.Role{
		ID:     role.ID,
		UID:    role.UID,
		Label:  role.Label,
		System: role.System,
		Permissions: sharedUtils.Map(role.Permissions, func(p dbModel.PermissionEntity) dllModel.Permission {
			return dllModel.Permission{
				UID:   p.UID,
				Label: p.Label,
			}
		}),
	}
}
