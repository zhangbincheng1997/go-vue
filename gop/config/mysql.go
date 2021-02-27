package config

// MySQL struct
type MySQL struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Protocol string `mapstructure:"protocol" json:"protocol" yaml:"protocol"`
	Database string `mapstructure:"database" json:"database" yaml:"database"`
	Config   string `mapstructure:"config" json:"config" yaml:"config"`
}
