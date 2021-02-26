package config

// JWT struct
type JWT struct {
	Key     string `yaml:"key" json:"key"`
	Timeout int    `yaml:"timeout" json:"timeout"`
}
