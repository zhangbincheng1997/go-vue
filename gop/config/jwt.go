package config

// JWT struct
type JWT struct {
	Realm         string `yaml:"realm" json:"realm"`
	Key           string `yaml:"key" json:"key"`
	Timeout       int    `yaml:"timeout" json:"timeout"`
	TokenLookup   string `yaml:"token-lookup" json:"tokenLookup"`
	TokenHeadName string `yaml:"token-head-name" json:"tokenHeadName"`
}
