package dll2db

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
)

func ToDBModelAPIKey(k dllModel.APIKey) dbModel.APIKeyEntity {
	model := dbModel.APIKeyEntity{
		ID:             k.ID.GetPayloadOrDefault(0),
		Label:          k.Label,
		ExpiresAt:      k.ExpiresAt,
		Revoked:        k.Revoked,
		RateLimit:      k.RateLimit,
		IPRestrictions: ToDBModelAPIKeyIPRestrictions(k),
	}
	if k.KeyHash != nil {
		model.KeyHash = *k.KeyHash
	}
	return model
}

func ToDBModelAPIKeyIPRestrictions(k dllModel.APIKey) []dbModel.APIKeyIPRestrictionEntity {
	return mapIPRestrictions(k.ID.GetPayloadOrDefault(0), k.IPRestrictions)
}

func mapIPRestrictions(apiKeyID uint32, cidrs []string) []dbModel.APIKeyIPRestrictionEntity {
	result := make([]dbModel.APIKeyIPRestrictionEntity, 0, len(cidrs))
	for _, cidr := range cidrs {
		result = append(result, dbModel.APIKeyIPRestrictionEntity{
			APIKeyID: apiKeyID,
			CIDR:     cidr,
		})
	}
	return result
}
