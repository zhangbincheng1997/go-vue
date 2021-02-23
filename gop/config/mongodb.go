package config

// MongoDB struct
type MongoDB struct {
	Addr     string `yaml:"addr" json:"addr"`
	Database string `yaml:"database" json:"database"`
}
