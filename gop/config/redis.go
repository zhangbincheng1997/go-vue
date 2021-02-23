package config

// Redis struct
type Redis struct {
	DB       int    `yaml:"db" json:"db"`
	Addr     string `yaml:"addr" json:"addr"`
	Password string `yaml:"password" json:"password"`
}
