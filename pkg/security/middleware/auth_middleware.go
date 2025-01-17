package middleware

import (
	"context"
	"github.com/GolangSpring/gospring/pkg/security/service"
	"net/http"
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

		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing token, please login first", http.StatusUnauthorized)
			return
		}

		// Token is valid - set custom headers or context (e.g., user info)
		// Add user info to the request context
		userClaims, err := middleware.AuthService.ParseUserClaims(token.Value)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", userClaims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
