package config

// Config struct
type Config struct {
	Zap     Zap     `yaml:"zap" json:"zap"`
	JWT     JWT     `yaml:"jwt" json:"jwt"`
	MySQL   MySQL   `yaml:"mysql" json:"mysql"`
	Redis   Redis   `yaml:"redis" json:"redis"`
	MongoDB MongoDB `yaml:"mongodb" json:"mongodb"`
	Qiniu   Qiniu   `yaml:"qiniu" json:"qiniu"`
}
