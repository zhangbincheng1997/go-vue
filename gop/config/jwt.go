package config

// JWT struct
type JWT struct {
	Key     string `mapstructure:"key" json:"key" yaml:"key"`
	Timeout int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
}
