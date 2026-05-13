/**
 * @file sdInstanceGroups.go
 * @brief Mapování vstupů skupin instancí zdrojů dat z GraphQL modelu do doménového modelu.
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
package gql2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelSDInstanceGroup(sdInstanceGroupInput graphQLModel.SDInstanceGroupInput) dllModel.SDInstanceGroup {
	return dllModel.SDInstanceGroup{
		ID:             sharedUtils.NewEmptyOptional[uint32](),
		Label:          sdInstanceGroupInput.Label,
		UserIdentifier: sdInstanceGroupInput.UserIdentifier,
		SDInstanceIDs:  sdInstanceGroupInput.SdInstanceIDs,
	}
}
