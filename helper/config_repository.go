package helper

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigRepository[T any] struct {
	FilePath string
}

func NewConfigRepository[T any](filePath string) *ConfigRepository[T] {
	return &ConfigRepository[T]{
		FilePath: filePath,
	}
}

func (repository *ConfigRepository[T]) Read() (*T, error) {
	fileContent, err := os.ReadFile(repository.FilePath)
	if err != nil {
		return nil, err
	}
	var config T
	if err := yaml.Unmarshal(fileContent, &config); err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (repository *ConfigRepository[T]) Save(config *T) error {
	file, err := os.Create(repository.FilePath)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Warn().Msgf("Failed to close file: %v", err)
		}
	}(file)

	_yamlOut, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if _, err = file.Write(_yamlOut); err != nil {
		return err
	}
	return nil
}
