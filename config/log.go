package config

import "github.com/gofiber/fiber/v2/log"

type LogConfig struct {
	Level log.Level
}

func NewLogConfig() *LogConfig {
	level := getString("LOG_LEVEL", "info")

	var logLevel log.Level
	switch level {
	case "trace":
		logLevel = log.LevelTrace
	case "debug":
		logLevel = log.LevelDebug
	case "info":
		logLevel = log.LevelInfo
	case "warn":
		logLevel = log.LevelWarn
	case "error":
		logLevel = log.LevelError
	default:
		logLevel = log.LevelInfo
	}

	return &LogConfig{Level: logLevel}
}
