package controller

import (
	"fmt"
	"github.com/GolangSpring/gospring/application"
	"github.com/GolangSpring/gospring/helper"
	"github.com/GolangSpring/gospring/pkg/security/middleware"
	"github.com/GolangSpring/gospring/pkg/security/repository"
	"github.com/GolangSpring/gospring/pkg/security/service"
	"github.com/go-fuego/fuego"
	"net/http"
	"time"
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
	fuego.Post(server, "/api-admin/assign-roles", controller.AssignRoles)
	fuego.Get(server, "/api-admin/all-roles", controller.AllRoles)

	fuego.Get(server, "/api-private/current-user", controller.CurrentUser)
	fuego.Get(server, "/api-private/logout", controller.Logout)

	fuego.Post(server, "/api-public/login", controller.Login)
	fuego.Post(server, "/api-public/register", controller.RegisterUser)
}

func (controller *AuthController) Health(c fuego.ContextNoBody) (string, error) {

	return fmt.Sprintf("%s, ok"), nil
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

	duration := time.Duration(24) * time.Hour
	helper.WriteTokenCookie(c, token, duration)

	return &http.Response{StatusCode: http.StatusNoContent, Status: token}, nil
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

func (controller *AuthController) Logout(c fuego.ContextNoBody) (*http.Response, error) {
	helper.WriteTokenCookie(c, "", -1)
	return &http.Response{StatusCode: http.StatusNoContent}, nil
}

func (controller *AuthController) CurrentUser(c fuego.ContextNoBody) (*service.UserClaims, error) {
	user, ok := c.Request().Context().Value("user").(*service.UserClaims)
	if !ok {
		return nil, fuego.HTTPError{
			Detail: "User not found within context",
			Status: http.StatusUnauthorized,
		}
	}

	return user, nil

}

func (controller *AuthController) AllRoles(c fuego.ContextNoBody) ([]string, error) {
	return controller.CasbinService.GetAllUsedRoles()
}
