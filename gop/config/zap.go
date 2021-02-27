package config

// Zap struct
type Zap struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	Filename   string `mapstructure:"filename" json:"filename" yaml:"filename"`
	MaxSize    int    `mapstructure:"max-size" json:"maxSize" yaml:"max-size"`
	MaxBackups int    `mapstructure:"max-backups" json:"maxBackups" yaml:"max-backups"`
	MaxAge     int    `mapstructure:"max-age" json:"maxAge" yaml:"max-age"`
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}
