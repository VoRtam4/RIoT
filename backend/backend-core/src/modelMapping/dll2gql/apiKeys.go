package dll2gql

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelAPIKey(k dllModel.APIKey) graphQLModel.APIKey {
	var expiresAt *string
	if k.ExpiresAt != nil {
		s := k.ExpiresAt.Format(time.RFC3339)
		expiresAt = &s
	}
	return graphQLModel.APIKey{
		ID:             k.ID.GetPayload(),
		Label:          k.Label,
		ExpiresAt:      expiresAt,
		Revoked:        k.Revoked,
		RateLimit:      k.RateLimit,
		Permissions:    k.Permissions,
		IPRestrictions: k.IPRestrictions,
	}
}
