package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelRole(role dbModel.RoleEntity) dllModel.Role {
	return dllModel.Role{
		ID:    role.ID,
		Label: role.Label,
		Permissions: sharedUtils.Map(role.Permissions, func(p dbModel.PermissionEntity) string {
			return p.Label
		}),
	}
}
