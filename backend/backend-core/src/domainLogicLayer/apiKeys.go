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

func UpdateAPIKeyForUser(userID uint32, id uint32, input graphQLModel.APIKeyInput) error {
	db := dbClient.GetRelationalDatabaseClientInstance()
	load := db.LoadAPIKeyByID(id)
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
	updated.ID = sharedUtils.NewOptionalOf(id)
	result := db.UpdateAPIKey(updated)
	if result.IsFailure() {
		return result.GetError()
	}
	return nil
}

func DeleteAPIKeyForUser(userID uint32, id uint32) error {
	db := dbClient.GetRelationalDatabaseClientInstance()
	load := db.LoadAPIKeyByID(id)
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
	return db.DeleteAPIKey(id)
}
