package domainLogicLayer

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func LoadUserRole(userID uint32) sharedUtils.Result[string] {
	return dbClient.GetRelationalDatabaseClientInstance().LoadUserRole(userID)
}

func AssignRoleToUser(userID uint32, roleID uint32) sharedUtils.Result[bool] {
	err := dbClient.GetRelationalDatabaseClientInstance().AssignRoleToUser(userID, roleID)
	if err != nil {
		return sharedUtils.NewFailureResult[bool](err)
	}
	return sharedUtils.NewSuccessResult(true)
}

func LoadPermissionsForUser(userID uint32) sharedUtils.Result[[]dbModel.PermissionEntity] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadPermissionsForUser(userID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]dbModel.PermissionEntity](result.GetError())
	}
	return sharedUtils.NewSuccessResult(result.GetPayload())
}

func LoadRoleByLabel(label string) sharedUtils.Result[dbModel.RoleEntity] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByLabel(label)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[dbModel.RoleEntity](result.GetError())
	}
	return sharedUtils.NewSuccessResult[dbModel.RoleEntity](result.GetPayload())
}
