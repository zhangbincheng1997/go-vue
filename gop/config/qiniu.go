package config

// Qiniu ...
type Qiniu struct {
	AccessKey     string `json:"accessKey" yaml:"access-key"`
	SecretKey     string `json:"secretKey" yaml:"secret-key"`
	ImgPath       string `json:"imgPath" yaml:"img-path"`
	Bucket        string `json:"bucket" yaml:"bucket"`
	Zone          string `json:"zone" yaml:"zone"`
	UseHTTPS      bool   `json:"useHttps" yaml:"use-https"`
	UseCdnDomains bool   `json:"useCdnDomains" yaml:"use-cdn-domains"`
}
