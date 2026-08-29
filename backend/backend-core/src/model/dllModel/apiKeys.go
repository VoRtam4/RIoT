/**
 * @file apiKeys.go
 * @brief Doménový model API klíčů a jejich omezení.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dllModel

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type APIKey struct {
	ID             sharedUtils.Optional[uint32]
	UID            string
	UserID         *uint32
	KeyHash        *string
	Label          string
	ExpiresAt      *time.Time
	Revoked        bool
	RateLimit      *uint32
	LastUsedAt     *time.Time
	Permissions    []string
	IPRestrictions []string
}
