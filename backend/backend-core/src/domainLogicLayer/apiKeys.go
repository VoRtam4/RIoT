package domainLogicLayer

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/db/dbClient"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func LoadAPIKeyByHash(hash string) sharedUtils.Result[sharedUtils.Optional[dbModel.APIKeyEntity]] {
	return dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeyByHash(hash)
}

func LoadAPIKeyByID(id uint32) sharedUtils.Result[sharedUtils.Optional[dbModel.APIKeyEntity]] {
	return dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeyByID(id)
}

func CreateAPIKey(userID uint32, roleID uint32, label string, expiresAt *time.Time) sharedUtils.Result[string] {
	rawKey := sharedUtils.GenerateRandomAlphanumericString(32)
	hash := sharedUtils.GenerateHexHash(rawKey)
	apiKey := dbModel.APIKeyEntity{
		UserID:    userID,
		RoleID:    roleID,
		KeyHash:   hash,
		Label:     label,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	result := dbClient.GetRelationalDatabaseClientInstance().PersistAPIKey(apiKey)
	if result.IsFailure() {
		return sharedUtils.NewFailureResult[string](result.GetError())
	}
	return sharedUtils.NewSuccessResult(rawKey)
}

func LoadAPIKeysForUser(userID uint32) sharedUtils.Result[[]dbModel.APIKeyEntity] {
	return dbClient.GetRelationalDatabaseClientInstance().LoadAPIKeysForUser(userID)
}

func UpdateAPIKey(apiKey dbModel.APIKeyEntity) error {
	return dbClient.GetRelationalDatabaseClientInstance().UpdateAPIKey(apiKey)
}

func DeleteAPIKey(id uint32) error {
	return dbClient.GetRelationalDatabaseClientInstance().DeleteAPIKey(id)
}
