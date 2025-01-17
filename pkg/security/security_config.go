package security

import "github.com/GolangSpring/gospring/application"

type SecurityConfig struct {
	Security struct {
		Secret string `yaml:"secret" validate:"required"`
	} `yaml:"security" validate:"required"`
}

func MustNewSecurityConfig(configPath string) *SecurityConfig {
	return application.MustNewConfigFromFile[SecurityConfig](configPath)
}
