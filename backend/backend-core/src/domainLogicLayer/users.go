/**
 * @file users.go
 * @brief Doménová logika pro uživatelskou konfiguraci a správu uživatelů.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace funkcionality souboru.
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
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func LoadUsers() sharedUtils.Result[[]graphQLModel.User] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadUsers()
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.User](result.GetError())
	}
	users := make([]graphQLModel.User, 0, len(result.GetPayload()))
	for _, user := range result.GetPayload() {
		mapped, err := mapUserWithRole(user)
		if err != nil {
			return sharedUtils.NewFailureResult[[]graphQLModel.User](err)
		}
		users = append(users, mapped)
	}
	return sharedUtils.NewSuccessResult(users)
}

func LoadUserByUID(uid string) sharedUtils.Result[graphQLModel.User] {
	user, err := loadUserDLLByUID(uid)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	mapped, err := mapUserWithRole(user)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	return sharedUtils.NewSuccessResult(mapped)
}

func LoadUserByID(id uint32) sharedUtils.Result[graphQLModel.User] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadUser(uint(id))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.User](result.GetError())
	}
	mapped, err := mapUserWithRole(result.GetPayload())
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	return sharedUtils.NewSuccessResult(mapped)
}

func UpdateUser(uid string, input graphQLModel.UserUpdateInput) sharedUtils.Result[graphQLModel.User] {
	user, err := loadUserDLLByUID(uid)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Name != nil {
		user.Name = optionalStringFromUpdate(*input.Name)
	}
	if input.ProfileImageURL != nil {
		user.ProfileImageURL = optionalStringFromUpdate(*input.ProfileImageURL)
	}
	result := dbClient.GetRelationalDatabaseClientInstance().UpdateUser(user)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.User](result.GetError())
	}
	mapped, err := mapUserWithRole(result.GetPayload())
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	return sharedUtils.NewSuccessResult(mapped)
}

func DisableUser(actorUserID uint32, uid string, reason *string) sharedUtils.Result[graphQLModel.User] {
	user, err := loadUserDLLByUID(uid)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	if user.ID.IsPresent() && uint32(user.ID.GetPayload()) == actorUserID {
		return sharedUtils.NewFailureResult[graphQLModel.User](fmt.Errorf("cannot disable current user"))
	}
	result := dbClient.GetRelationalDatabaseClientInstance().SetUserDisabled(user.ID.GetPayload(), true, reason)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.User](result.GetError())
	}
	if err := dbClient.GetRelationalDatabaseClientInstance().RevokeUserSessions(user.ID.GetPayload()); err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	mapped, err := mapUserWithRole(result.GetPayload())
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	return sharedUtils.NewSuccessResult(mapped)
}

func EnableUser(uid string) sharedUtils.Result[graphQLModel.User] {
	user, err := loadUserDLLByUID(uid)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	result := dbClient.GetRelationalDatabaseClientInstance().SetUserDisabled(user.ID.GetPayload(), false, nil)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.User](result.GetError())
	}
	mapped, err := mapUserWithRole(result.GetPayload())
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.User](err)
	}
	return sharedUtils.NewSuccessResult(mapped)
}

func RevokeUserSessions(uid string) error {
	user, err := loadUserDLLByUID(uid)
	if err != nil {
		return err
	}
	return dbClient.GetRelationalDatabaseClientInstance().RevokeUserSessions(user.ID.GetPayload())
}

func LoadUserSessions(userID uint32) sharedUtils.Result[[]graphQLModel.UserSession] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadUserSessions(uint(userID))
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.UserSession](result.GetError())
	}
	userResult := dbClient.GetRelationalDatabaseClientInstance().LoadUser(uint(userID))
	if userResult.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.UserSession](userResult.GetError())
	}
	return sharedUtils.NewSuccessResult(mapUserSessions(result.GetPayload(), userResult.GetPayload().UID))
}

func LoadUserSessionsByUserUID(userUID string) sharedUtils.Result[[]graphQLModel.UserSession] {
	user, err := loadUserDLLByUID(userUID)
	if err != nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.UserSession](err)
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadUserSessions(user.ID.GetPayload())
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.UserSession](result.GetError())
	}
	return sharedUtils.NewSuccessResult(mapUserSessions(result.GetPayload(), user.UID))
}

func RevokeSession(actorUserID uint32, sessionUID string, allowForeign bool) error {
	normalizedUID, normalizeErr := normalizeUserSessionUID(sessionUID)
	if normalizeErr != nil {
		return normalizeErr
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadUserSessionByUID(normalizedUID)
	if result.IsFailure() {
		return result.GetError()
	}
	sessionOptional := result.GetPayload()
	if sessionOptional.IsEmpty() {
		return fmt.Errorf("not found")
	}
	session := sessionOptional.GetPayload()
	if uint32(session.UserID) != actorUserID && !allowForeign {
		return fmt.Errorf("forbidden")
	}
	return dbClient.GetRelationalDatabaseClientInstance().RevokeUserSessionByUID(normalizedUID)
}

func RevokeOwnOtherSessions(userID uint32, currentSessionID uint) error {
	return dbClient.GetRelationalDatabaseClientInstance().RevokeOtherUserSessions(uint(userID), currentSessionID)
}

func GetUserConfig(id uint32) sharedUtils.Result[graphQLModel.UserConfig] {
	loadResult := dbClient.GetRelationalDatabaseClientInstance().LoadUserConfig(id)
	if loadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.UserConfig](loadResult.GetError())
	}
	userConfig := loadResult.GetPayload()
	userLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadUser(uint(id))
	if userLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.UserConfig](userLoadResult.GetError())
	}
	userConfig.UserUID = userLoadResult.GetPayload().UID
	return sharedUtils.NewSuccessResult[graphQLModel.UserConfig](dll2gql.ToGraphQLModelUserConfig(userConfig))
}

func UpdateUserConfig(id uint32, input graphQLModel.UserConfigInput) sharedUtils.Result[graphQLModel.UserConfig] {
	userConfig := gql2dll.ToDLLModelUserConfig(input)
	userConfig.UserID = id
	if sdInstanceGroupPersistResult := dbClient.GetRelationalDatabaseClientInstance().PersistUserConfig(userConfig); sdInstanceGroupPersistResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.UserConfig](sdInstanceGroupPersistResult.GetError())
	}
	userLoadResult := dbClient.GetRelationalDatabaseClientInstance().LoadUser(uint(id))
	if userLoadResult.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.UserConfig](userLoadResult.GetError())
	}
	userConfig.UserUID = userLoadResult.GetPayload().UID
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelUserConfig(userConfig))
}

func DeleteUserConfig(id uint32) error {
	return dbClient.GetRelationalDatabaseClientInstance().DeleteUserConfig(id)
}

func loadUserDLLByUID(uid string) (dllModel.User, error) {
	normalizedUID, normalizeErr := normalizeUserUID(uid)
	if normalizeErr != nil {
		return dllModel.User{}, normalizeErr
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadUserBasedOnUID(normalizedUID)
	if result.IsFailure() {
		return dllModel.User{}, result.GetError()
	}
	userOptional := result.GetPayload()
	if userOptional.IsEmpty() || userOptional.GetPayload().ID.IsEmpty() {
		return dllModel.User{}, fmt.Errorf("not found")
	}
	return userOptional.GetPayload(), nil
}

func mapUserWithRole(user dllModel.User) (graphQLModel.User, error) {
	if user.ID.IsEmpty() {
		return graphQLModel.User{}, fmt.Errorf("missing id")
	}
	roleResult := dbClient.GetRelationalDatabaseClientInstance().GetUserRole(uint32(user.ID.GetPayload()))
	if roleResult.IsFailure() {
		return graphQLModel.User{}, roleResult.GetError()
	}
	role := roleResult.GetPayload()
	return dll2gql.ToGraphQLModelUser(user, &role), nil
}

func mapUserSessions(sessions []dllModel.UserSession, userUID string) []graphQLModel.UserSession {
	return sharedUtils.Map(sessions, func(session dllModel.UserSession) graphQLModel.UserSession {
		return dll2gql.ToGraphQLModelUserSession(session, userUID)
	})
}

func optionalStringFromUpdate(value string) sharedUtils.Optional[string] {
	if value == "" {
		return sharedUtils.NewEmptyOptional[string]()
	}
	return sharedUtils.NewOptionalOf(value)
}
