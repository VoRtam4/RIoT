package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDBModelRole(role dllModel.Role) dbModel.RoleEntity {
	return dbModel.RoleEntity{
		ID:    role.ID,
		Label: role.Label,
	}
}
