package middleware

import (
	"context"
	"fmt"
	"github.com/GolangSpring/gospring/helper"
	"github.com/GolangSpring/gospring/pkg/security/service"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	EmptyString    = ""
	UserContextKey = "user"
)

type AuthMiddleware struct {
	AuthService   service.IAuthService
	CasbinService *service.CasbinService
}

func (middleware *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj := r.URL.Path // The resource (API endpoint)
		act := r.Method   // The HTTP method (GET, POST, etc.)

		isPublicAllowed, err := middleware.CasbinService.HasPermission(service.CasbinPublicKey, obj, act)
		if err == nil && isPublicAllowed {
			next.ServeHTTP(w, r)
			return
		}

		var tokenFound string
		authToken := r.Header.Get("Authorization")
		if authToken != EmptyString {
			parts := strings.SplitN(authToken, " ", 2)
			var bearerPrefix, token string
			if err := helper.SliceUnpack(parts, &bearerPrefix, &token); err != nil {
				errString := fmt.Sprintf("Invalid Authorization header: %v", err)
				log.Warn().Msg(errString)
				http.Error(w, errString, http.StatusBadRequest)
				return
			}

			if strings.EqualFold(bearerPrefix, "Bearer") {
				tokenFound = token
			} else {
				log.Warn().Msg("Authorization header does not contain 'Bearer' prefix")
				http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
				return
			}
		}

		// Fallback to cookie if no Authorization header
		if tokenFound == EmptyString {
			if token, err := r.Cookie("token"); err == nil {
				tokenFound = token.Value
			}
		}

		if tokenFound == EmptyString {
			http.Error(w, "Missing token, unauthorized.", http.StatusUnauthorized)
			return
		}

		// Token is valid - set custom headers or context (e.g., user info)
		// Add user info to the request context
		userClaims, err := middleware.AuthService.ParseUserClaims(tokenFound)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, userClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
