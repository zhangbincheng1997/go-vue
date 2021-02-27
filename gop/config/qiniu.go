package config

// Qiniu ...
type Qiniu struct {
	AccessKey     string `mapstructure:"access-key" json:"accessKey" yaml:"access-key"`
	SecretKey     string `mapstructure:"secret-key" json:"secretKey" yaml:"secret-key"`
	ImgPath       string `mapstructure:"img-path" json:"imgPath" yaml:"img-path"`
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Zone          string `mapstructure:"zone" json:"zone" yaml:"zone"`
	UseHTTPS      bool   `mapstructure:"use-https" json:"useHttps" yaml:"use-https"`
	UseCdnDomains bool   `mapstructure:"use-cdn-domains" json:"useCdnDomains" yaml:"use-cdn-domains"`
}
