package domainLogicLayer

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func LoadRolesByLabels(labels []string) sharedUtils.Result[[]graphQLModel.Role] {
	roles := make([]graphQLModel.Role, 0, len(labels))
	for _, label := range labels {
		result := dbClient.GetRelationalDatabaseClientInstance().GetRoleIDByLabel(label)
		if result.IsFailure() {
			return sharedUtils.NewFailureResult[[]graphQLModel.Role](result.GetError())
		}
		role := result.GetPayload()
		roles = append(roles, dll2gql.ToGraphQLModelRole(role))
	}
	return sharedUtils.NewSuccessResult(roles)
}

func LoadUserRole(userID uint32) sharedUtils.Result[graphQLModel.Role] {
	result := dbClient.GetRelationalDatabaseClientInstance().GetUserRole(userID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelRole(result.GetPayload()))
}

func AssignRoleToUser(userID uint32, roleID uint32) error {
	result := dbClient.GetRelationalDatabaseClientInstance().GetUserRole(userID)
	if result.IsFailure() {
		return result.GetError()
	}
	return dbClient.GetRelationalDatabaseClientInstance().SetUserRole(userID, roleID)
}
