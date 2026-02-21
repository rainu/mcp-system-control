package model

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/rainu/go-yacl"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type DebugConfig struct {
	LogLevel       string      `yaml:"log-level,omitempty"`
	LogLevelParsed *slog.Level `yaml:"-"`
}

func (d *DebugConfig) SetDefaults() {
	if d.LogLevel == "" {
		d.LogLevel = "info"
		d.LogLevelParsed = yacl.P(slog.LevelInfo)
	}
}

func (d *DebugConfig) GetUsage(field string) string {
	switch field {
	case "LogLevel":
		return fmt.Sprintf("Log level (%s, %s, %s, %s)", LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError)
	}
	return ""
}

func (d *DebugConfig) Validate() error {
	switch strings.ToLower(d.LogLevel) {
	case LogLevelDebug:
		d.LogLevelParsed = yacl.P(slog.LevelDebug)
	case LogLevelInfo:
		d.LogLevelParsed = yacl.P(slog.LevelInfo)
	case LogLevelWarn:
		d.LogLevelParsed = yacl.P(slog.LevelWarn)
	case LogLevelError:
		d.LogLevelParsed = yacl.P(slog.LevelError)
	default:
		return fmt.Errorf("Invalid log level '%s'", d.LogLevel)
	}

	return nil
}
