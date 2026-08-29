/**
 * @file apiKeys.go
 * @brief Mapování API klíčů z doménového modelu do GraphQL modelu.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dll2gql

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelAPIKey(k dllModel.APIKey) graphQLModel.APIKey {
	var expiresAt *string
	if k.ExpiresAt != nil {
		s := k.ExpiresAt.Format(time.RFC3339Nano)
		expiresAt = &s
	}
	var lastUsedAt *string
	if k.LastUsedAt != nil {
		s := k.LastUsedAt.Format(time.RFC3339Nano)
		lastUsedAt = &s
	}
	return graphQLModel.APIKey{
		UID:            k.UID,
		Label:          k.Label,
		ExpiresAt:      expiresAt,
		Revoked:        k.Revoked,
		RateLimit:      k.RateLimit,
		LastUsedAt:     lastUsedAt,
		Permissions:    k.Permissions,
		IPRestrictions: k.IPRestrictions,
	}
}
