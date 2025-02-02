package security

import (
	"github.com/GolangSpring/gospring/application"
	"github.com/GolangSpring/gospring/pkg/security/service"
)

type SecurityConfig struct {
	Security struct {
		Secret string `yaml:"secret" validate:"required"`
	} `yaml:"security" validate:"required"`
	Smtp *service.SmtpConfig `yaml:"smtp" validate:"required"`
}

func MustNewSecurityConfig(configPath string) *SecurityConfig {
	return application.MustNewConfigFromFile[SecurityConfig](configPath)
}
