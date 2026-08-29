/**
 * @file apiKeys.go
 * @brief Doménová logika správy API klíčů, oprávnění a IP restrikcí.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality správy API klíčů v doménové vrstvě.
 *
 * @ingroup riot_backend_core
 */
package domainLogicLayer

import (
	"fmt"
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/dll2gql"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/modelMapping/gql2dll"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func LoadAPIKeyByHash(hash string) sharedUtils.Result[sharedUtils.Optional[dllModel.APIKey]] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeyByHash(hash)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.APIKey]](result.GetError())
	}
	apiKeyOpt := result.GetPayload()
	if apiKeyOpt.IsEmpty() {
		return sharedUtils.NewSuccessResult(sharedUtils.NewEmptyOptional[dllModel.APIKey]())
	}
	apiKey := apiKeyOpt.GetPayload()
	if apiKey.Revoked {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.APIKey]](fmt.Errorf("api key revoked"))
	}
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return sharedUtils.NewFailureResult[sharedUtils.Optional[dllModel.APIKey]](fmt.Errorf("api key expired"))
	}
	return sharedUtils.NewSuccessResult(sharedUtils.NewOptionalOf(apiKey))
}

func LoadAPIKeyByUID(userID uint32, uid string) sharedUtils.Result[graphQLModel.APIKey] {
	normalizedUID, normalizeErr := normalizeAPIKeyUID(uid)
	if normalizeErr != nil {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](normalizeErr)
	}
	result := dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeyByUID(normalizedUID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](result.GetError())
	}
	if result.GetPayload().IsEmpty() {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](fmt.Errorf("not found"))
	}
	r := result.GetPayload().GetPayload()
	if userID != *r.UserID {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](fmt.Errorf("not found"))
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelAPIKey(r))
}

func LoadAPIKeyByID(userID uint32, id uint32) sharedUtils.Result[graphQLModel.APIKey] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeyByID(id)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](result.GetError())
	}
	if result.GetPayload().IsEmpty() {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](fmt.Errorf("not found"))
	}
	r := result.GetPayload().GetPayload()
	if userID != *r.UserID {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](fmt.Errorf("not found"))
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelAPIKey(r))
}

func LoadAPIKeysForUser(userID uint32) sharedUtils.Result[[]graphQLModel.APIKey] {
	result := dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeysForUser(userID)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[[]graphQLModel.APIKey](result.GetError())
	}
	return sharedUtils.NewSuccessResult(sharedUtils.Map(result.GetPayload(), func(k dllModel.APIKey) graphQLModel.APIKey {
		return dll2gql.ToGraphQLModelAPIKey(k)
	}))
}

func LoadAPIKeysByUserUID(userUID string) sharedUtils.Result[[]graphQLModel.APIKey] {
	user, err := loadUserDLLByUID(userUID)
	if err != nil {
		return sharedUtils.NewFailureResult[[]graphQLModel.APIKey](err)
	}
	return LoadAPIKeysForUser(uint32(user.ID.GetPayload()))
}

func CreateAPIKey(userID uint32, input graphQLModel.APIKeyInput) sharedUtils.Result[string] {
	if input.Permissions != nil {
		roleResult := LoadUserRole(userID)
		if roleResult.IsFailure() {
			return sharedUtils.NewFailureResult[string](roleResult.GetError())
		}
		userRole := roleResult.GetPayload()
		err := sharedUtils.ValidatePermissions[graphQLModel.Permission](userRole.Permissions, input.Permissions, func(p graphQLModel.Permission) string { return p.UID })
		if err != nil {
			return sharedUtils.NewFailureResult[string](err)
		}
	}
	rawKey := sharedUtils.GenerateRandomAlphanumericString(32)
	hash := sharedUtils.GenerateHexHash(rawKey)
	k := gql2dll.ToDLLModelAPIKey(input)
	k.KeyHash = &hash
	k.Revoked = false
	result := dbClient.GetRelationalDatabaseClientInstance().CreateAPIKey(userID, k)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[string](result.GetError())
	}
	return sharedUtils.NewSuccessResult(rawKey)
}

func UpdateAPIKeyForUser(userID uint32, uid string, input graphQLModel.APIKeyInput) error {
	db := dbClient.GetRelationalDatabaseClientInstance()
	normalizedUID, normalizeErr := normalizeAPIKeyUID(uid)
	if normalizeErr != nil {
		return normalizeErr
	}
	load := db.LoadAPIKeyByUID(normalizedUID)
	if load.IsFailure() {
		return load.GetError()
	}
	if load.GetPayload().IsEmpty() {
		return fmt.Errorf("not found")
	}
	k := load.GetPayload().GetPayload()
	if userID != *k.UserID {
		return fmt.Errorf("not found")
	}
	if input.Permissions != nil {
		roleResult := LoadUserRole(userID)
		if roleResult.IsFailure() {
			return roleResult.GetError()
		}
		userRole := roleResult.GetPayload()
		err := sharedUtils.ValidatePermissions[graphQLModel.Permission](userRole.Permissions, input.Permissions, func(p graphQLModel.Permission) string { return p.UID })
		if err != nil {
			return err
		}
	}
	updated := gql2dll.ApplyAPIKeyUpdate(k, input)
	if k.ID.IsEmpty() {
		return fmt.Errorf("missing id")
	}
	updated.ID = k.ID
	result := db.UpdateAPIKey(updated)
	if result.IsFailure() {
		return result.GetError()
	}
	return nil
}

func RevokeAPIKey(userID uint32, uid string, allowForeign bool) sharedUtils.Result[graphQLModel.APIKey] {
	k, err := loadAPIKeyForOperation(userID, uid, allowForeign)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](err)
	}
	k.Revoked = true
	result := dbClient.GetRelationalDatabaseClientInstance().UpdateAPIKey(k)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelAPIKey(result.GetPayload()))
}

func RotateAPIKey(userID uint32, uid string, allowForeign bool) sharedUtils.Result[string] {
	k, err := loadAPIKeyForOperation(userID, uid, allowForeign)
	if err != nil {
		return sharedUtils.NewFailureResult[string](err)
	}
	rawKey := sharedUtils.GenerateRandomAlphanumericString(32)
	hash := sharedUtils.GenerateHexHash(rawKey)
	result := dbClient.GetRelationalDatabaseClientInstance().UpdateAPIKeyHash(k.ID.GetPayload(), hash)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[string](result.GetError())
	}
	return sharedUtils.NewSuccessResult(rawKey)
}

func UpdateAPIKeyPermissions(userID uint32, uid string, permissionUIDs []string, allowForeign bool) sharedUtils.Result[graphQLModel.APIKey] {
	k, err := loadAPIKeyForOperation(userID, uid, allowForeign)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](err)
	}
	if k.UserID == nil {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](fmt.Errorf("missing user id"))
	}
	if err := validateAPIKeyPermissionsForUser(*k.UserID, permissionUIDs); err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](err)
	}
	k.Permissions = permissionUIDs
	result := dbClient.GetRelationalDatabaseClientInstance().UpdateAPIKey(k)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelAPIKey(result.GetPayload()))
}

func UpdateAPIKeyRestrictions(userID uint32, uid string, input graphQLModel.APIKeyRestrictionsInput, allowForeign bool) sharedUtils.Result[graphQLModel.APIKey] {
	k, err := loadAPIKeyForOperation(userID, uid, allowForeign)
	if err != nil {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](err)
	}
	if input.ExpiresAt != nil {
		if *input.ExpiresAt == "" {
			k.ExpiresAt = nil
		} else if t, err := time.Parse(time.RFC3339, *input.ExpiresAt); err == nil {
			k.ExpiresAt = &t
		} else {
			return sharedUtils.NewFailureResult[graphQLModel.APIKey](err)
		}
	}
	if input.Revoked != nil {
		k.Revoked = *input.Revoked
	}
	if input.RateLimit != nil {
		val := *input.RateLimit
		k.RateLimit = &val
	}
	if input.IPRestrictions != nil {
		k.IPRestrictions = input.IPRestrictions
	}
	result := dbClient.GetRelationalDatabaseClientInstance().UpdateAPIKey(k)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[graphQLModel.APIKey](result.GetError())
	}
	return sharedUtils.NewSuccessResult(dll2gql.ToGraphQLModelAPIKey(result.GetPayload()))
}

func TouchAPIKeyLastUsedAt(id uint32) error {
	return dbClient.GetRelationalDatabaseClientInstance().TouchAPIKeyLastUsedAt(id)
}

func DeleteAPIKeyForUser(userID uint32, uid string) error {
	db := dbClient.GetRelationalDatabaseClientInstance()
	normalizedUID, normalizeErr := normalizeAPIKeyUID(uid)
	if normalizeErr != nil {
		return normalizeErr
	}
	load := db.LoadAPIKeyByUID(normalizedUID)
	if load.IsFailure() {
		return load.GetError()
	}
	if load.GetPayload().IsEmpty() {
		return fmt.Errorf("not found")
	}
	k := load.GetPayload().GetPayload()
	if userID != *k.UserID {
		return fmt.Errorf("not found")
	}
	if k.ID.IsEmpty() {
		return fmt.Errorf("missing id")
	}
	return db.DeleteAPIKey(k.ID.GetPayload())
}

func loadAPIKeyForOperation(userID uint32, uid string, allowForeign bool) (dllModel.APIKey, error) {
	normalizedUID, normalizeErr := normalizeAPIKeyUID(uid)
	if normalizeErr != nil {
		return dllModel.APIKey{}, normalizeErr
	}
	load := dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeyByUID(normalizedUID)
	if load.IsFailure() {
		return dllModel.APIKey{}, load.GetError()
	}
	if load.GetPayload().IsEmpty() {
		return dllModel.APIKey{}, fmt.Errorf("not found")
	}
	k := load.GetPayload().GetPayload()
	if k.UserID == nil || k.ID.IsEmpty() {
		return dllModel.APIKey{}, fmt.Errorf("missing id")
	}
	if userID != *k.UserID && !allowForeign {
		return dllModel.APIKey{}, fmt.Errorf("not found")
	}
	return k, nil
}

func validateAPIKeyPermissionsForUser(userID uint32, permissions []string) error {
	roleResult := LoadUserRole(userID)
	if roleResult.IsFailure() {
		return roleResult.GetError()
	}
	userRole := roleResult.GetPayload()
	return sharedUtils.ValidatePermissions[graphQLModel.Permission](userRole.Permissions, permissions, func(p graphQLModel.Permission) string { return p.UID })
}
