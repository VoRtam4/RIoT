/**
 * @file sdInstanceGroups.go
 * @brief Mapování skupin instancí zdrojů dat z doménového modelu do GraphQL modelu.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní mapování skupin instancí zdrojů dat.
 * - Vojtěch Hubáček: doplnění mapování labelů.
 *
 * @ingroup riot_backend_core
 */
package dll2gql

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelSDInstanceGroup(sdInstanceGroup dllModel.SDInstanceGroup) graphQLModel.SDInstanceGroup {
	return graphQLModel.SDInstanceGroup{
		UID:            sdInstanceGroup.UID,
		Label:          sdInstanceGroup.Label,
		UserIdentifier: sdInstanceGroup.UserIdentifier,
		SdInstanceUIDs: sdInstanceGroup.SDInstanceUIDs,
	}
}
