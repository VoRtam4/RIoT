/**
 * @file users.go
 * @brief Mapování uživatelské konfigurace z doménového modelu do GraphQL modelu.
 *
 * @author Michal Bureš
 *
 * @par Autorský podíl
 * - Michal Bureš: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_backend_core
 */
package dll2gql

import (
	"time"

	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/dllModel"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/model/graphQLModel"
)

func ToGraphQLModelUser(user dllModel.User, role *dllModel.Role) graphQLModel.User {
	var lastLoginAt *string
	if user.LastLoginAt.IsPresent() {
		s := user.LastLoginAt.GetPayload().Format(time.RFC3339Nano)
		lastLoginAt = &s
	}
	var disabledAt *string
	if user.DisabledAt.IsPresent() {
		s := user.DisabledAt.GetPayload().Format(time.RFC3339Nano)
		disabledAt = &s
	}
	var gqlRole *graphQLModel.Role
	if role != nil {
		mapped := ToGraphQLModelRole(*role)
		gqlRole = &mapped
	}
	return graphQLModel.User{
		UID:             user.UID,
		Username:        user.Username,
		Email:           user.Email,
		Name:            user.Name.ToPointer(),
		ProfileImageURL: user.ProfileImageURL.ToPointer(),
		Oauth2Provider:  user.OAuth2Provider.ToPointer(),
		LastLoginAt:     lastLoginAt,
		Disabled:        user.Disabled,
		DisabledAt:      disabledAt,
		DisabledReason:  user.DisabledReason.ToPointer(),
		Role:            gqlRole,
	}
}

func ToGraphQLModelUserSession(session dllModel.UserSession, userUID string) graphQLModel.UserSession {
	return graphQLModel.UserSession{
		UID:       session.UID,
		UserUID:   userUID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339Nano),
		Revoked:   session.Revoked,
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		CreatedAt: session.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt: session.UpdatedAt.Format(time.RFC3339Nano),
	}
}

func ToGraphQLModelUserConfig(userConfig dllModel.UserConfig) graphQLModel.UserConfig {
	return graphQLModel.UserConfig{
		UserUID: userConfig.UserUID,
		Config:  userConfig.Config,
	}
}
