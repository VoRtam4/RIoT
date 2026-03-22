package gql2dll

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToDLLModelAPIKey(input graphQLModel.APIKeyInput) dllModel.APIKey {
	var k dllModel.APIKey
	if input.Label != nil {
		k.Label = *input.Label
	}
	if input.ExpiresAt != nil {
		if t, err := time.Parse(time.RFC3339, *input.ExpiresAt); err == nil {
			k.ExpiresAt = &t
		}
	}
	if input.Revoked != nil {
		k.Revoked = *input.Revoked
	}
	if input.RateLimit != nil {
		val := *input.RateLimit
		k.RateLimit = &val
	}
	if input.Permissions != nil {
		k.Permissions = input.Permissions
	}
	if input.IPRestrictions != nil {
		k.IPRestrictions = input.IPRestrictions
	}
	return k
}

func ApplyAPIKeyUpdate(prev dllModel.APIKey, input graphQLModel.APIKeyInput) dllModel.APIKey {
	out := prev
	if input.Label != nil {
		out.Label = *input.Label
	}
	if input.ExpiresAt != nil {
		if *input.ExpiresAt == "" {
			out.ExpiresAt = nil
		} else if t, err := time.Parse(time.RFC3339, *input.ExpiresAt); err == nil {
			out.ExpiresAt = &t
		}
	}
	if input.Revoked != nil {
		out.Revoked = *input.Revoked
	}
	if input.RateLimit != nil {
		val := *input.RateLimit
		out.RateLimit = &val
	}
	if input.Permissions != nil {
		out.Permissions = input.Permissions
	}
	if input.IPRestrictions != nil {
		out.IPRestrictions = input.IPRestrictions
	}
	return out
}
