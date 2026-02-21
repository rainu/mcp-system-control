package approval

import (
	"log/slog"
	"mcp-system-control/approval/message"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

var defaultFormatter *message.Formatter

func init() {
	var err error
	defaultFormatter, err = message.NewFormatterAuto()
	if err != nil {
		slog.Warn("Failed to create auto formatter, using English", "error", err)
		defaultFormatter, _ = message.NewFormatter(message.LanguageEnglish)
	}
}

// SetLanguage sets the language for approval messages
func SetLanguage(lang string) error {
	var msgLang message.Language

	switch strings.ToLower(lang) {
	case "auto", "":
		msgLang = message.DetectLanguage()
	case "en", "english":
		msgLang = message.LanguageEnglish
	case "de", "german", "deutsch":
		msgLang = message.LanguageGerman
	default:
		msgLang = message.LanguageEnglish
	}

	formatter, err := message.NewFormatter(msgLang)
	if err != nil {
		return err
	}

	defaultFormatter = formatter
	return nil
}

// formatApprovalMessage formats the tool request into a human-readable message
func formatApprovalMessage(request *mcp.CallToolRequest) string {
	return defaultFormatter.Format(request)
}
