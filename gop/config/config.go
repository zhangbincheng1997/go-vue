package config

// Config struct
type Config struct {
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	MySQL   MySQL   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	MongoDB MongoDB `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
	Qiniu   Qiniu   `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
}
