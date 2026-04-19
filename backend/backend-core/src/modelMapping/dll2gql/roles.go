package dll2gql

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToGraphQLModelRole(role dllModel.Role) graphQLModel.Role {
	return graphQLModel.Role{
		ID:    role.ID,
		Label: role.Label,
		Permissions: sharedUtils.Map(role.Permissions, func(p dllModel.Permission) graphQLModel.Permission {
			return graphQLModel.Permission{
				UID:   p.UID,
				Label: p.Label,
			}
		}),
	}
}
