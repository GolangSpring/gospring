package application

import (
	"encoding/json"
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
	return MustNewConfigFromFile[Config](configPath)
}
