package service

import "errors"

var (
	UserNotFound        = errors.New("UserNotFound")
	UserExists          = errors.New("UserExists")
	UserAlreadyVerified = errors.New("UserAlreadyVerified")

	OtpIncorrect = errors.New("OtpIncorrect")
	OtpNotFound  = errors.New("OtpNotFound")

	TokenInvalid = errors.New("TokenInvalid")
	TokenExpired = errors.New("TokenExpired")

	ResetPasswordNotMatched = errors.New("ResetPasswordNotMatched")
)
