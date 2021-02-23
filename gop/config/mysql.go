package config

// MySQL struct
type MySQL struct {
	Addr     string `yaml:"addr" json:"addr"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Protocol string `yaml:"protocol" json:"protocol"`
	Database string `yaml:"database" json:"database"`
	Config   string `yaml:"config" json:"config"`
}
