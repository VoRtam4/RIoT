/**
 * @file setup.go
 * @brief Inicializace GraphQL handleru, transportů a autentizace GraphQL požadavků.
 *
 * @author Michal Bureš
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Michal Bureš: původní nastavení GraphQL serveru.
 * - Vojtěch Hubáček: napojení na společný server, autorizace operací a API key přístup v GraphQL transportech.
 *
 * @ingroup riot_backend_core
 */
package graphql

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/api/graphql/gsc"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/auth"
	"github.com/MichalBures-OG/bp-bures-RIoT-backend-core/src/domainLogicLayer"
	"github.com/MichalBures-OG/bp-bures-RIoT-commons/src/sharedUtils"
	"github.com/gorilla/websocket"
)

func GetHandler() http.Handler {
	allowedOrigins := sharedUtils.NewSetFromSlice(strings.Split(sharedUtils.GetEnvironmentVariableValue("ALLOWED_ORIGINS").GetPayloadOrDefault("http://localhost:8080,http://localhost:1234"), ","))
	graphQLServer := handler.New(gsc.NewExecutableSchema(gsc.Config{Resolvers: new(Resolver)}))
	graphQLServer.AddTransport(transport.POST{})
	graphQLServer.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				return allowedOrigins.Contains(origin)
			},
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
			if _, ok := auth.PrincipalFromContext(ctx); ok {
				return ctx, nil, nil
			}

			apiKeyRaw := strings.TrimSpace(initPayload.GetString("X-API-Key"))
			if apiKeyRaw == "" {
				apiKeyRaw = strings.TrimSpace(initPayload.GetString("x-api-key"))
				if apiKeyRaw == "" {
					apiKeyRaw = strings.TrimSpace(initPayload.GetString("X-API-KEY"))
				}
			}
			if apiKeyRaw != "" {
				hash := sharedUtils.GenerateHexHash(apiKeyRaw)
				apiKeyResult := domainLogicLayer.LoadAPIKeyByHash(hash)
				if apiKeyResult.IsFailure() {
					return ctx, nil, fmt.Errorf("unauthorized")
				}
				apiKeyOpt := apiKeyResult.GetPayload()
				if apiKeyOpt.IsEmpty() {
					return ctx, nil, fmt.Errorf("unauthorized")
				}
				apiKey := apiKeyOpt.GetPayload()
				if apiKey.UserID == nil || apiKey.ID.IsEmpty() {
					return ctx, nil, fmt.Errorf("invalid api key")
				}
				roleResult := domainLogicLayer.LoadUserRole(*apiKey.UserID)
				if roleResult.IsFailure() {
					return ctx, nil, fmt.Errorf("unauthorized")
				}
				rolePermissions := make(map[string]bool, len(roleResult.GetPayload().Permissions))
				for _, permission := range roleResult.GetPayload().Permissions {
					rolePermissions[permission.UID] = true
				}
				allowed := make(map[string]bool, len(apiKey.Permissions))
				for _, permission := range apiKey.Permissions {
					if rolePermissions[permission] {
						allowed[permission] = true
					}
				}
				apiKeyID := apiKey.ID.GetPayload()
				principal := auth.Principal{
					Type:              auth.PrincipalAPIKey,
					UserID:            *apiKey.UserID,
					APIKeyID:          &apiKeyID,
					AllowedOperations: allowed,
					SynchronizedAt:    time.Now(),
				}
				return auth.ContextWithPrincipal(ctx, &principal), nil, nil
			}

			tokenStr := strings.TrimPrefix(initPayload.Authorization(), "Bearer ")
			if tokenStr == "" {
				return ctx, nil, fmt.Errorf("missing auth")
			}
			token, err := auth.ParseJWT(tokenStr)
			if err != nil {
				return ctx, nil, fmt.Errorf("unauthorized")
			}
			if !auth.IsJWTValid(token) {
				return ctx, nil, fmt.Errorf("unauthorized")
			}
			subject, err := token.Claims.GetSubject()
			if err != nil {
				return ctx, nil, fmt.Errorf("invalid subject")
			}
			userIDUint64, err := strconv.ParseUint(subject, 10, 32)
			if err != nil {
				return ctx, nil, fmt.Errorf("invalid user id")
			}
			principal := auth.Principal{
				Type:   auth.PrincipalUserSession,
				UserID: uint32(userIDUint64),
			}
			return auth.ContextWithPrincipal(ctx, &principal), nil, nil
		},
	})
	graphQLServer.Use(extension.Introspection{})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		bodyStr := string(bodyBytes)
		if strings.Contains(bodyStr, "__schema") {
			fmt.Println("[GRAPHQL] introspection → skipping auth")
			graphQLServer.ServeHTTP(w, r)
			return
		}
		auth.JWTAuthenticationMiddleware(graphQLServer).ServeHTTP(w, r)
	})
}

func authorizeOperation(ctx context.Context, operation string, opType string) (*auth.Principal, error) {
	principal, ok := auth.PrincipalFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	if !auth.CanAccessOperation(principal, operation, opType) {
		return nil, fmt.Errorf("forbidden")
	}
	return principal, nil
}
