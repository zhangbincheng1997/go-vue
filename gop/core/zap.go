package core

import (
	"encoding/json"

	"go.uber.org/zap"
)

// Zap ...
func Zap() *zap.Logger {
	rawJSON := []byte(`{
		"level": "info",
		"outputPaths": ["stdout", "logs/log.log"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return logger
}
