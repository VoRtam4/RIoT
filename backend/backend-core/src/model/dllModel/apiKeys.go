package dllModel

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type APIKey struct {
	ID             sharedUtils.Optional[uint32]
	UserID         *uint32
	KeyHash        *string
	Label          string
	ExpiresAt      *time.Time
	Revoked        bool
	RateLimit      *uint32
	Permissions    []string
	IPRestrictions []string
}
