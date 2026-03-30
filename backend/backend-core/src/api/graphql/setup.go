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
			tokenRaw, ok := initPayload["Authorization"]
			if !ok {
				return ctx, nil, fmt.Errorf("missing auth")
			}
			tokenStr, ok := tokenRaw.(string)
			if !ok {
				return ctx, nil, fmt.Errorf("invalid auth format")
			}
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
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
