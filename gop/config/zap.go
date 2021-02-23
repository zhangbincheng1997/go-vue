package config

// Zap struct
type Zap struct {
	Level         string `json:"level" yaml:"level"`
	Format        string `json:"format" yaml:"format"`
	Prefix        string `json:"prefix" yaml:"prefix"`
	Director      string `json:"director"  yaml:"director"`
	LinkName      string `json:"linkName" yaml:"link-name"`
	ShowLine      bool   `json:"showLine" yaml:"showLine"`
	EncodeLevel   string `json:"encodeLevel" yaml:"encode-level"`
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktrace-key"`
	LogInConsole  bool   `json:"logInConsole" yaml:"log-in-console"`
}
