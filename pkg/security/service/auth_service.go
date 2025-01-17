package service

import (
	"context"
	"fmt"
	. "github.com/GolangSpring/gospring/pkg/security/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserClaims struct {
	ID                 uint     `json:"id"`
	UserName           string   `json:"user_name"`
	Roles              []string `json:"roles"`
	ExpirationDuration float64  `json:"exp"`
	IsVerified         bool     `json:"is_verified"`
}

func NewUserClaims(userID uint, userName string, roles []string, expiration float64, isVerified bool) *UserClaims {
	return &UserClaims{
		ID:                 userID,
		UserName:           userName,
		Roles:              roles,
		ExpirationDuration: expiration,
		IsVerified:         isVerified,
	}
}

func (claims *UserClaims) Validate() error {
	if claims.ExpirationDuration < float64(time.Now().Unix()) {
		return jwt.ErrTokenExpired
	}
	return nil
}

func NewUser(name string, email string, password string) (*User, error) {
	user := User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &user, nil
}

type IAuthService interface {
	LoginWithEmail(ctx context.Context, email string, password string) (string, error)
	LoginWithUserName(ctx context.Context, userName string, password string) (string, error)
	IssueJsonWebToken(claims *jwt.MapClaims) string
	IssueLoginToken(user *User, expiration time.Duration) (string, error)
	ExtractUserClaims(claims *jwt.MapClaims) (*UserClaims, error)

	ParseUserClaims(tokenString string) (*UserClaims, error)
	RegisterUser(ctx context.Context, name string, email string, password string) (*User, error)
	AssignRoles(ctx context.Context, userID uint, roles []string) (*User, error)
}

var _ IAuthService = (*AuthService)(nil)

type AuthService struct {
	Secret      string
	UserService IUserService
}

func (service *AuthService) PostConstruct() {}

func (service *AuthService) AssignRoles(ctx context.Context, userID uint, roles []string) (*User, error) {
	return service.UserService.UpdateUserRoles(ctx, userID, roles)
}

func NewAuthService(userService IUserService, secret string) *AuthService {
	return &AuthService{
		UserService: userService,
		Secret:      secret,
	}
}

func (service *AuthService) RegisterUser(ctx context.Context, name string, email string, rawPassword string) (*User, error) {
	hashedPassword, err := service.GenerateHashedPassword(rawPassword)
	if err != nil {
		return nil, err
	}
	user, err := NewUser(name, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, service.UserService.AddUser(ctx, user)
}

func (service *AuthService) ExtractUserClaims(claims *jwt.MapClaims) (*UserClaims, error) {

	userID, ok := (*claims)["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'id' claim")
	}
	userName, ok := (*claims)["user_name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'user_name' claim")
	}

	roleInterfaces, ok := (*claims)["roles"].([]any)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'roles' claim")
	}
	roles := make([]string, len(roleInterfaces))
	for idx, role := range roleInterfaces {
		roles[idx] = role.(string)
	}

	expiration, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'exp' claim")
	}

	isVerified, ok := (*claims)["is_verified"].(bool)
	if !ok {
		return nil, fmt.Errorf("invalid or missing 'is_verified' claim")
	}

	userClaims := NewUserClaims(uint(userID), userName, roles, expiration, isVerified)
	return userClaims, nil
}

func (service *AuthService) DecodeJsonWebTokenWithSecret(rawToken string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (service *AuthService) DecodeJsonWebToken(rawToken string) (*jwt.Token, error) {
	return service.DecodeJsonWebTokenWithSecret(rawToken, []byte(service.Secret))

}

func (service *AuthService) LoginWithUserName(ctx context.Context, userName string, password string) (string, error) {
	user, err := service.UserService.FindByUserName(ctx, userName)
	if err != nil {
		return "", UserNotFound
	}

	if err := service.VerifyPassword(password, user.Password); err != nil {
		return "", err
	}

	return service.IssueLoginToken(user, time.Hour)
}

func (service *AuthService) LoginWithEmail(ctx context.Context, email string, password string) (string, error) {
	user, err := service.UserService.FindByEmail(ctx, email)
	if err != nil {
		return "", UserNotFound
	}

	if err := service.VerifyPassword(password, user.Password); err != nil {
		return "", err
	}

	return service.IssueLoginToken(user, time.Hour)
}

func (service *AuthService) IssueJsonWebToken(claims *jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(service.Secret))
	log.Info().Msgf("Issue Token: %v", tokenString)
	return tokenString
}

func (service *AuthService) IssueLoginToken(user *User, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_name":   user.Name,
		"id":          user.ID,
		"roles":       user.Roles,
		"exp":         time.Now().Add(expiration).Unix(),
		"is_verified": user.IsVerified,
	}
	return service.IssueJsonWebToken(&claims), nil
}

func (service *AuthService) ParseUserClaims(tokenString string) (*UserClaims, error) {
	_jwt, err := service.DecodeJsonWebToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := _jwt.Claims.(jwt.MapClaims); ok && _jwt.Valid {
		userClaims, err := service.ExtractUserClaims(&claims)
		if err != nil {
			return nil, err
		}

		return userClaims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

func (service *AuthService) GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (service *AuthService) VerifyPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
