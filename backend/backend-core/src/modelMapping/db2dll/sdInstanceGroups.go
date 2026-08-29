/**
 * @file sdInstanceGroups.go
 * @brief Mapování skupin instancí zdrojů dat z databázového modelu do doménového modelu.
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
package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelSDInstanceGroup(sdInstanceGroupEntity dbModel.SDInstanceGroupEntity) dllModel.SDInstanceGroup {
	return dllModel.SDInstanceGroup{
		ID:             sharedUtils.NewOptionalOf(sdInstanceGroupEntity.ID),
		UID:            sdInstanceGroupEntity.UID,
		Label:          sdInstanceGroupEntity.Label,
		UserIdentifier: sdInstanceGroupEntity.UserIdentifier,
		SDInstanceIDs: sharedUtils.Map(sdInstanceGroupEntity.GroupMembershipRecords, func(sdInstanceGroupMembershipEntity dbModel.SDInstanceGroupMembershipEntity) uint32 {
			return sdInstanceGroupMembershipEntity.SDInstanceID
		}),
		SDInstanceUIDs: sharedUtils.Map(sdInstanceGroupEntity.GroupMembershipRecords, func(sdInstanceGroupMembershipEntity dbModel.SDInstanceGroupMembershipEntity) string {
			return sdInstanceGroupMembershipEntity.SDInstance.UID
		}),
	}
}
