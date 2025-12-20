package config

import (
	"os"

	"github.com/rs/zerolog"
)

type LogConfig struct {
	Level  zerolog.Level
	Pretty bool
}

func NewLogConfig() *LogConfig {
	level := getString("LOG_LEVEL", "info")
	pretty := getBool("LOG_PRETTY", false)

	return &LogConfig{
		Level:  parseLogLevel(level),
		Pretty: pretty,
	}
}

func parseLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

// NewLogger creates a zerolog.Logger based on config
func (c *LogConfig) NewLogger() zerolog.Logger {
	var logger zerolog.Logger

	switch c.Pretty {
	case true:
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	default:
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(c.Level)

	return logger
}
