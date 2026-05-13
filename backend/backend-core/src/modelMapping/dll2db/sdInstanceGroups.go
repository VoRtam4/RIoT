/**
 * @file sdInstanceGroups.go
 * @brief Mapování skupin instancí zdrojů dat z doménového modelu do databázového modelu.
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
package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDBModelEntitySDInstanceGroup(sdInstanceGroup dllModel.SDInstanceGroup) dbModel.SDInstanceGroupEntity {
	sdInstanceGroupID := sdInstanceGroup.ID.GetPayloadOrDefault(0)
	return dbModel.SDInstanceGroupEntity{
		ID:             sdInstanceGroupID,
		Label:          sdInstanceGroup.Label,
		UserIdentifier: sdInstanceGroup.UserIdentifier,
		GroupMembershipRecords: sharedUtils.Map[uint32, dbModel.SDInstanceGroupMembershipEntity](sdInstanceGroup.SDInstanceIDs, func(sdInstanceID uint32) dbModel.SDInstanceGroupMembershipEntity {
			return dbModel.SDInstanceGroupMembershipEntity{
				SDInstanceGroupID: sdInstanceGroupID,
				SDInstanceID:      sdInstanceID,
			}
		}),
	}
}
