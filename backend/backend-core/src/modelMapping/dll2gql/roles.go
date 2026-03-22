package dll2gql

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelRole(role dllModel.Role) graphQLModel.Role {
	return graphQLModel.Role{
		ID:          role.ID,
		Label:       role.Label,
		Permissions: role.Permissions,
	}
}
