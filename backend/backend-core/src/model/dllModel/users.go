/**
 * @file users.go
 * @brief Doménový model uživatelů, session dat a uživatelské konfigurace.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní uživatelský model a session data.
 * - Vojtěch Hubáček: doplnění vazby uživatele na roli.
 *
 * @ingroup riot_backend_core
 */
package dllModel

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
)

type User struct {
	ID                     sharedUtils.Optional[uint]
	UID                    string
	RoleID                 uint32
	Username               string
	Email                  string
	Name                   sharedUtils.Optional[string]
	ProfileImageURL        sharedUtils.Optional[string]
	OAuth2Provider         sharedUtils.Optional[string]
	OAuth2ProviderIssuedID sharedUtils.Optional[string]
	LastLoginAt            sharedUtils.Optional[time.Time]
	Disabled               bool
	DisabledAt             sharedUtils.Optional[time.Time]
	DisabledReason         sharedUtils.Optional[string]
	Sessions               []UserSession
	// TODO: Implement 'Invocations', 'UserConfig' and other possibly missing fields as needed
}

type UserSession struct {
	ID               sharedUtils.Optional[uint] // TODO: Consider getting rid of Optional[T] within dllModel...
	UID              string
	UserID           uint
	RefreshTokenHash string
	ExpiresAt        time.Time
	Revoked          bool
	IPAddress        string
	UserAgent        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserConfig struct {
	UserID  uint32
	UserUID string
	Config  string
}
