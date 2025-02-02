package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/ugurcsen/gods-generic/sets/hashset"
	"text/template"
	"time"
)

type UserVerificationClaims struct {
	ID                 uint    `json:"id"`
	ExpirationDuration float64 `json:"exp"`
}

type UserVerificationService struct {
	SmtpService                       ISmtpService
	UserService                       IUserService
	AuthService                       IAuthService
	OtpService                        IOtpService
	AdminPushedEmailVerificationCache *hashset.Set[uint]
}

func NewUserVerificationService(
	smtpService ISmtpService,
	userService IUserService,
	authService IAuthService,
	otpService IOtpService,
) *UserVerificationService {
	return &UserVerificationService{
		SmtpService:                       smtpService,
		UserService:                       userService,
		AuthService:                       authService,
		OtpService:                        otpService,
		AdminPushedEmailVerificationCache: hashset.New[uint](),
	}
}

func (service *UserVerificationService) IssueVerificationToken(ctx context.Context, userID uint) (string, error) {
	user, err := service.UserService.FindByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.IsVerified {
		return "", UserAlreadyVerified
	}

	claims := jwt.MapClaims{
		"purpose": string(PurposeGuestEmailVerification),
		"id":      user.ID,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}

	_jwt := service.AuthService.IssueJsonWebToken(&claims)

	return _jwt, nil
}

func (service *UserVerificationService) IsAdminAskingForVerification(userID uint) bool {
	return service.AdminPushedEmailVerificationCache.Contains(userID)
}

func (service *UserVerificationService) SendVerificationEmailByUserID(ctx context.Context, userID uint, isAdminPushed bool) error {

	user, err := service.UserService.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.IsVerified {
		return UserAlreadyVerified
	}

	if isAdminPushed {
		log.Info().Msgf("Admin is pushing email verification for user %d", user.ID)
		service.AdminPushedEmailVerificationCache.Add(user.ID)
	}

	_template, err := template.New("email_verification").Parse(EMAIL_VERIFICATION_HTML_TEMPLATE)
	if err != nil {
		return err
	}
	otp := service.OtpService.GenerateOtp(user.ID, PurposeGuestEmailVerification)

	var buffer bytes.Buffer
	emailTemplate := NewEmailTemplate(user.Name, otp.Code, service.SmtpService.GetSmtpConfig().CompanyName)
	if err := _template.Execute(&buffer, emailTemplate); err != nil {
		return err
	}

	emailContent := buffer.String()

	message := service.SmtpService.CreateNewMessage(user.Email, "Email Verification", emailContent, ContentTypeHtml)

	return service.SmtpService.SendEmail(message)
}

func (service *UserVerificationService) SendVerificationEmailByToken(ctx context.Context, token string) error {
	claims, err := service.parseVerificationClaims(token)
	if err != nil {
		return err
	}
	return service.SendVerificationEmailByUserID(ctx, claims.ID, false)
}

func (service *UserVerificationService) parseVerificationClaims(token string) (*UserVerificationClaims, error) {
	_jwt, err := service.AuthService.DecodeJsonWebToken(token)
	if err != nil {
		return nil, err
	}
	if claims, ok := _jwt.Claims.(jwt.MapClaims); ok && _jwt.Valid {
		verificationClaims, err := service.extractVerificationClaims(&claims)
		if err != nil {
			return nil, err
		}
		return verificationClaims, nil
	}
	return nil, TokenInvalid
}

func (service *UserVerificationService) extractVerificationClaims(claims *jwt.MapClaims) (*UserVerificationClaims, error) {
	var verificationClaims UserVerificationClaims
	purpose, ok := (*claims)["purpose"].(string)
	if !ok || purpose != string(PurposeGuestEmailVerification) {
		return nil, fmt.Errorf("invalid or missing 'purpose' claim, getting %s, expects %v", purpose, PurposeGuestEmailVerification)
	}
	userID, ok := (*claims)["id"].(float64)
	if !ok {
		return nil, TokenInvalid
	}
	verificationClaims.ID = uint(userID)
	exp, ok := (*claims)["exp"].(float64)
	if !ok {
		return nil, TokenInvalid
	}
	verificationClaims.ExpirationDuration = exp
	return &verificationClaims, nil
}

func (service *UserVerificationService) VerifyEmail(token string, otpCode string) error {
	claims, err := service.parseVerificationClaims(token)
	if err != nil {
		return err
	}

	if err := service.OtpService.VerifyOtp(claims.ID, PurposeGuestEmailVerification, otpCode); err != nil {
		return err
	}
	user, err := service.UserService.FindByID(context.Background(), claims.ID)
	if err != nil {
		return err
	}

	if service.AdminPushedEmailVerificationCache.Contains(user.ID) {
		log.Info().Msgf("user: %d is inside cache, removing from cache upone verification", user.ID)
		service.AdminPushedEmailVerificationCache.Remove(user.ID)
	}

	return service.UserService.ActivateUser(context.Background(), user)
}
