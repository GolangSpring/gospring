package application

type LogMode string

const (
	LogModeJson LogMode = "json"
	LogModeText LogMode = "text"
)

// LogConfig defines the configuration for file rotation
type LogConfig struct {
	FileName   string  `yaml:"file_name" validate:"required"`
	MaxSize    int     `yaml:"max_size" validate:"required"`    // in MB
	MaxBackups int     `yaml:"max_backups" validate:"required"` // number of backups
	MaxAge     int     `yaml:"max_age" validate:"required"`     // in days
	Compress   bool    `yaml:"compress" validate:"required"`
	LogMode    LogMode `yaml:"log_mode" validate:"required"`
}
