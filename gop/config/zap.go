package config

// Zap struct
type Zap struct {
	Level      string `json:"level" yaml:"level"`
	Filename   string `json:"filename" yaml:"filename"`
	MaxSize    int    `json:"maxSize" yaml:"max-size"`
	MaxBackups int    `json:"maxBackups" yaml:"max-backups"`
	MaxAge     int    `json:"maxAge" yaml:"max-age"`
	Compress   bool   `json:"compress" yaml:"compress"`
}
