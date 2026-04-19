package db2dll

import (
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dbModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

func ToDLLModelAPIKey(e dbModel.APIKeyEntity) dllModel.APIKey {
	permissions := make([]string, 0)
	for _, p := range e.Permissions {
		permissions = append(permissions, p.UID)
	}
	return dllModel.APIKey{
		ID:          sharedUtils.NewOptionalOf(e.ID),
		UserID:      &e.UserID,
		KeyHash:     &e.KeyHash,
		Label:       e.Label,
		ExpiresAt:   e.ExpiresAt,
		Revoked:     e.Revoked,
		RateLimit:   e.RateLimit,
		Permissions: permissions,
		IPRestrictions: sharedUtils.Map(
			e.IPRestrictions,
			func(r dbModel.APIKeyIPRestrictionEntity) string {
				return r.CIDR
			},
		),
	}
}
