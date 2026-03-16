package rules

import (
	"strings"

	"github.com/n1ckerr0r/loglint/internal/log_message"
)

// Позволяет пользователю настроить слайс чувствительных данных
type SensitiveRule struct {
	keywords      []string
	caseSensitive bool
}

func NewSensitiveRule() *SensitiveRule {
	return &SensitiveRule{
		keywords:      getKeywords(),
		caseSensitive: false,
	}
}

func (sensitiveRule *SensitiveRule) CustomizeKeywords(newKeywords []string, caseSensitive bool) {
	sensitiveRule.keywords = make([]string, len(newKeywords))
	copy(sensitiveRule.keywords, newKeywords)

	sensitiveRule.caseSensitive = caseSensitive
}

func getKeywords() []string {

	return []string{
		"password",
		"token",
		"api_key",
		"apikey",
		"secret",
		"auth",
		"private_key",
		"access_token",
		"refresh_token",
		"jwt",
		"session",
		"cookie",
	}
}

func (r *SensitiveRule) Name() string {
	return "sensitive"
}

func (sensitiveRule *SensitiveRule) Check(msg log_message.LogMessage) error {

	text := msg.Text

	if !sensitiveRule.caseSensitive {
		text = strings.ToLower(text)
	}

	for _, kw := range sensitiveRule.keywords {

		keyword := kw

		if !sensitiveRule.caseSensitive {
			keyword = strings.ToLower(kw)
		}

		if strings.Contains(text, keyword) {
			return ErrSensitiveData
		}
	}

	return nil
}
