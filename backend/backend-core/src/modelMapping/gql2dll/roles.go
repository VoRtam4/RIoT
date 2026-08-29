/**
 * @file roles.go
 * @brief Mapování rolí a oprávnění z GraphQL modelu do doménového modelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package gql2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelRole(role graphQLModel.Role) dllModel.Role {
	return dllModel.Role{
		UID:   role.UID,
		Label: role.Label,
		Permissions: sharedUtils.Map(role.Permissions, func(p graphQLModel.Permission) dllModel.Permission {
			return dllModel.Permission{
				UID:   p.UID,
				Label: p.Label,
			}
		}),
	}
}
