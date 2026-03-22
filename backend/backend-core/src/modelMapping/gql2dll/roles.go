package gql2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToDLLModelRole(role graphQLModel.Role) dllModel.Role {
	return dllModel.Role{
		ID:          role.ID,
		Label:       role.Label,
		Permissions: role.Permissions,
	}
}
