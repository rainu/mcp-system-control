package message

import (
	"os"
	"strings"
)

// Language represents a supported language
type Language string

const (
	LanguageEnglish Language = "en"
	LanguageGerman  Language = "de"
)

// DetectLanguage detects the system language
func DetectLanguage() Language {
	// Try LANG environment variable first
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = os.Getenv("LANGUAGE")
	}
	if lang == "" {
		lang = os.Getenv("LC_ALL")
	}
	if lang == "" {
		lang = os.Getenv("LC_MESSAGES")
	}

	// Parse language code (e.g., "de_DE.UTF-8" -> "de")
	lang = strings.ToLower(lang)
	if strings.HasPrefix(lang, "de") {
		return LanguageGerman
	}

	// Default to English
	return LanguageEnglish
}
