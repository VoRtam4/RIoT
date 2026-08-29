/**
 * @file roles.go
 * @brief Doménová logika pro načítání uživatelských rolí a oprávnění.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality rolí v doménové vrstvě.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func LoadRoles() sharedUtils.Result[[]graphQLModel.Role] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadRoles()
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), dll2gql.ToGraphQLModelRole))
}

func LoadPermissions() sharedUtils.Result[[]graphQLModel.Permission] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadPermissions()
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.Permission](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), func(p dllModel.Permission) graphQLModel.Permission {
		return graphQLModel.Permission{UID: p.UID, Label: p.Label}
	}))
}

func LoadRoleByUID(uid string) sharedUtils.Result[graphQLModel.Role] {
	normalizedUID, normalizeErr := normalizeRoleUID(uid)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.Role](normalizeErr)
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByUID(normalizedUID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelRole(result.GetPayload()))
}

func LoadRolesByUIDs(uids []string) sharedUtils.Result[[]graphQLModel.Role] {
	roles := make([]graphQLModel.Role, 0, len(uids))
	for _, uid := range uids {
		normalizedUID, normalizeErr := normalizeRoleUID(uid)
		if normalizeErr != nil {
			return sharedUtils.NewFailureResult[[]graphQLModel.Role](normalizeErr)
		}
		result := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByUID(normalizedUID)
		if result.IsFailure() {
			return sharedUtils.NewFailureResult[[]graphQLModel.Role](result.GetError())
		}
		role := result.GetPayload()
		roles = append(roles, dll2gql.ToGraphQLModelRole(role))
	}
	return sharedUtils.NewSuccessResult(roles)
}

func CreateRole(input graphQLModel.RoleInput) sharedUtils.Result[graphQLModel.Role] {
	role := dllModel.Role{
		Label:       input.Label,
		System:      false,
		Permissions: permissionUIDsToDLL(input.PermissionUIDs),
	}
	result := dbClient.GetRelationalDatabaseClientInstance().CreateRole(role)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelRole(result.GetPayload()))
}

func UpdateRole(uid string, input graphQLModel.RoleInput) sharedUtils.Result[graphQLModel.Role] {
	role, err := loadRoleDLLByUID(uid)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.Role](err)
	}
	role.Label = input.Label
	role.Permissions = permissionUIDsToDLL(input.PermissionUIDs)
	result := dbClient.GetRelationalDatabaseClientInstance().UpdateRole(role)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelRole(result.GetPayload()))
}

func DeleteRole(uid string) error {
	role, err := loadRoleDLLByUID(uid)
	if err != nil {
		return err
	}
	if role.System {
		return fmt.Errorf("cannot delete system role")
	}
	countResult := dbClient.GetRelationalDatabaseClientInstance().CountUsersWithRole(role.ID)
	if countResult.IsFailure() {
		return countResult.GetError()
	}
	if countResult.GetPayload() > 0 {
		return fmt.Errorf("cannot delete role assigned to users")
	}
	return dbClient.GetRelationalDatabaseClientInstance().DeleteRole(role.ID)
}

func CloneRole(uid string, label string) sharedUtils.Result[graphQLModel.Role] {
	role, err := loadRoleDLLByUID(uid)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.Role](err)
	}
	role.ID = 0
	role.UID = ""
	role.Label = label
	role.System = false
	result := dbClient.GetRelationalDatabaseClientInstance().CreateRole(role)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelRole(result.GetPayload()))
}

func UpdatePermissionLabel(uid string, label string) sharedUtils.Result[graphQLModel.Permission] {
	result := dbClient.GetRelationalDatabaseClientInstance().UpdatePermissionLabel(uid, label)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Permission](result.GetError())
	}
	permission := result.GetPayload()
	return sharedUtils.NewSuccessResult(graphQLModel.Permission{UID: permission.UID, Label: permission.Label})
}

func LoadRoleByLabel(label string) sharedUtils.Result[graphQLModel.Role] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByLabel(label)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelRole(result.GetPayload()))
}

func LoadRoleIDByLabel(label string) sharedUtils.Result[uint32] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByLabel(label)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[uint32](result.GetError())
	}
	return sharedUtils.NewSuccessResult(result.GetPayload().ID)
}

func LoadRoleIDByUID(uid string) sharedUtils.Result[uint32] {
	normalizedUID, normalizeErr := normalizeRoleUID(uid)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[uint32](normalizeErr)
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByUID(normalizedUID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[uint32](result.GetError())
	}
	return sharedUtils.NewSuccessResult(result.GetPayload().ID)
}

func LoadUserRole(userID uint32) sharedUtils.Result[graphQLModel.Role] {
	result := dbClient.GetRelationalDatabaseClientInstance().GetUserRole(userID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelRole(result.GetPayload()))
}

func LoadUserRoleByUID(userUID string) sharedUtils.Result[graphQLModel.Role] {
	normalizedUserUID, normalizeErr := normalizeUserUID(userUID)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.Role](normalizeErr)
	}
	userResult := dbClient.GetRelationalDatabaseClientInstance().LoadUserBasedOnUID(normalizedUserUID)
	if userResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](userResult.GetError())
	}
	userOptional := userResult.GetPayload()
	if userOptional.IsEmpty() || userOptional.GetPayload().ID.IsEmpty() {
		return sharedUtils.NewFailureResult[graphQLModel.Role](fmt.Errorf("not found"))
	}
	return LoadUserRole(uint32(userOptional.GetPayload().ID.GetPayload()))
}

func AssignRoleToUser(userUID string, roleUID string) error {
	normalizedUserUID, normalizeErr := normalizeUserUID(userUID)
	if normalizeErr != nil {
		return normalizeErr
	}
	normalizedRoleUID, normalizeErr := normalizeRoleUID(roleUID)
	if normalizeErr != nil {
		return normalizeErr
	}
	userResult := dbClient.GetRelationalDatabaseClientInstance().LoadUserBasedOnUID(normalizedUserUID)
	if userResult.IsFailure() {
		return userResult.GetError()
	}
	userOptional := userResult.GetPayload()
	if userOptional.IsEmpty() || userOptional.GetPayload().ID.IsEmpty() {
		return fmt.Errorf("not found")
	}
	userID := uint32(userOptional.GetPayload().ID.GetPayload())
	result := dbClient.GetRelationalDatabaseClientInstance().GetUserRole(userID)
	if result.IsFailure() {
		return result.GetError()
	}
	roleResult := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByUID(normalizedRoleUID)
	if roleResult.IsFailure() {
		return roleResult.GetError()
	}
	return dbClient.GetRelationalDatabaseClientInstance().SetUserRole(userID, roleResult.GetPayload().ID)
}

func loadRoleDLLByUID(uid string) (dllModel.Role, error) {
	normalizedUID, normalizeErr := normalizeRoleUID(uid)
	if normalizeErr != nil {
		return dllModel.Role{}, normalizeErr
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadRoleByUID(normalizedUID)
	if result.IsFailure() {
		return dllModel.Role{}, result.GetError()
	}
	return result.GetPayload(), nil
}

func permissionUIDsToDLL(permissionUIDs []string) []dllModel.Permission {
	return sharedUtils.Map(permissionUIDs, func(uid string) dllModel.Permission {
		return dllModel.Permission{UID: uid}
	})
}
