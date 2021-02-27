package config

// MongoDB struct
type MongoDB struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Database string `mapstructure:"database" json:"database" yaml:"database"`
}
