package controller

import (
	"github.com/go-fuego/fuego"
	"go-spring/application"
	"go-spring/pkg/security/middleware"
	"go-spring/pkg/security/service"
	"net/http"
)

var _ application.IController = (*CasbinController)(nil)

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
	casbinMiddleware := middleware.CasbinMiddleware{
		CasbinService: controller.CasbinService,
		AuthService:   controller.AuthService,
	}
	return []func(next http.Handler) http.Handler{
		casbinMiddleware.Middleware,
	}
}
