package application

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

func MustNewConfigFromFile[T any](configPath string) *T {
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal().Msgf("Failed to read config file: %v", err)
	}
	var config T
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatal().Msgf("Failed to unmarshal config: %v", err)
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		log.Fatal().Msgf("Failed to validate config: %v", err)
	}
	return &config
}
