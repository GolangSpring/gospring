package controller

import (
	"github.com/go-fuego/fuego"
	"go-spring/application"
	"go-spring/pkg/security/middleware"
	"go-spring/pkg/security/repository"
	"go-spring/pkg/security/service"
	"net/http"
)

type LoginBody struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type UserCredentials struct {
	UserName string `json:"user_name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RoleAssignBody struct {
	UserID uint     `json:"user_id" validate:"required"`
	Roles  []string `json:"roles" validate:"required"`
}

var _ application.IController = (*AuthController)(nil)

type AuthController struct {
	AuthService   service.IAuthService
	CasbinService *service.CasbinService
}

func NewAuthController(authService service.IAuthService, casbinService *service.CasbinService) *AuthController {
	return &AuthController{
		AuthService:   authService,
		CasbinService: casbinService,
	}
}

func (controller *AuthController) Routes(server *fuego.Server) {
	fuego.Post(server, "/api-public/login", controller.Login)
	fuego.Post(server, "/api-public/register", controller.RegisterUser)
	fuego.Post(server, "/api-admin/assign-roles", controller.AssignRoles)
}

func (controller *AuthController) Login(c fuego.ContextWithBody[LoginBody]) (*http.Response, error) {
	loginBody, err := c.Body()
	if err != nil {
		return nil, err
	}
	var loginErr error
	var token string
	if len(loginBody.Email) != 0 {
		token, loginErr = controller.AuthService.LoginWithEmail(c.Request().Context(), loginBody.Email, loginBody.Password)
	}
	if len(loginBody.UserName) != 0 {
		token, loginErr = controller.AuthService.LoginWithUserName(c.Request().Context(), loginBody.UserName, loginBody.Password)
	}

	if loginErr != nil {
		return nil, fuego.HTTPError{
			Detail: loginErr.Error(),
			Status: http.StatusBadRequest,
		}
	}
	responseCookie := http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
	c.SetCookie(responseCookie)

	return &http.Response{StatusCode: http.StatusNoContent}, nil
}

func (controller *AuthController) RegisterUser(c fuego.ContextWithBody[UserCredentials]) (any, error) {
	credentials, err := c.Body()
	if err != nil {
		return nil, err
	}

	user, err := controller.AuthService.RegisterUser(c.Request().Context(), credentials.UserName, credentials.Email, credentials.Password)
	if err != nil {
		return nil, fuego.HTTPError{
			Detail: err.Error(),
			Status: http.StatusBadRequest,
		}
	}
	return user, nil
}

func (controller *AuthController) Middlewares() []func(next http.Handler) http.Handler {

	authMiddle := middleware.AuthMiddleware{
		AuthService:   controller.AuthService,
		CasbinService: controller.CasbinService,
	}

	return []func(next http.Handler) http.Handler{
		authMiddle.Middleware,
	}
}

func (controller *AuthController) AssignRoles(c fuego.ContextWithBody[RoleAssignBody]) (*repository.User, error) {
	rolesBody, err := c.Body()
	if err != nil {
		return nil, err
	}
	user, err := controller.AuthService.AssignRoles(c.Request().Context(), rolesBody.UserID, rolesBody.Roles)
	if err != nil {
		return nil, fuego.HTTPError{
			Detail: err.Error(),
			Status: http.StatusBadRequest,
		}
	}
	return user, nil
}
