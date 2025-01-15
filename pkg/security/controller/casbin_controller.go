package controller

import (
	"github.com/go-fuego/fuego"
	"go-spring/application"
	"go-spring/pkg/security/service"
	"net/http"
)

var _ application.IController = (*CasbinController)(nil)

type CasbinMiddleware struct {
	AuthService   service.IAuthService
	CasbinService *service.CasbinService
}

func (middleware *CasbinMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj := r.URL.Path // The resource (API endpoint)
		act := r.Method   // The HTTP method (GET, POST, etc.)

		isPublicAllowed, err := middleware.CasbinService.HasPermission(service.CasbinPublicKey, obj, act)
		if err == nil && isPublicAllowed {
			next.ServeHTTP(w, r)
			return
		}

		userClaims, ok := r.Context().Value("user").(*service.UserClaims) // User info from context (e.g., JWT claims)
		if userClaims == nil || !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user is allowed to access the resource
		for _, role := range userClaims.Roles {
			allowed, err := middleware.CasbinService.Enforcer.Enforce(role, obj, act)
			if err != nil {
				http.Error(w, "Authorization error", http.StatusInternalServerError)
				return
			}
			if allowed {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}

type CasbinController struct {
	AuthService   service.IAuthService
	CasbinService *service.CasbinService
}

func NewCasbinController(casbinService *service.CasbinService, authService service.IAuthService) *CasbinController {
	return &CasbinController{
		AuthService:   authService,
		CasbinService: casbinService,
	}
}

func (controller *CasbinController) Routes(server *fuego.Server) {}

func (controller *CasbinController) Middlewares() []func(next http.Handler) http.Handler {
	casbinMiddleware := CasbinMiddleware{
		CasbinService: controller.CasbinService,
		AuthService:   controller.AuthService,
	}
	return []func(next http.Handler) http.Handler{
		casbinMiddleware.Middleware,
	}
}
