package application

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerConfig *ServerConfig `yaml:"server" validate:"required"`
	LogConfig    *LogConfig    `yaml:"log" validate:"required"`
}

func (config *Config) AsJson() string {
	_json, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		return ""
	}
	return string(_json)
}

func MustNewAppConfig(configPath string) *Config {
	config := MustNewConfigFromFile[Config](configPath)
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(config)
	if err != nil {
		log.Fatal().Msgf("Failed to validate config: %v", err)
	}
	return config
}
