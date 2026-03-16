package rules

import (
	"unicode"
	"unicode/utf8"

	"github.com/n1ckerr0r/loglint/internal/log_message"
)

// Это правило считает пустой лог или лог, начинающийся с заглавной буквы некорректным
type LowercaseRule struct{}

func NewLowercaseRule() *LowercaseRule {
	return &LowercaseRule{}
}

func (lowercaseRule *LowercaseRule) Name() string {
	return "LowercaseFirstRune"
}

func (lowercaseRule *LowercaseRule) Check(message log_message.LogMessage) error {
	if message.Text == "" {
		return nil
	}

	firstRune, _ := utf8.DecodeRuneInString(message.Text)

	if unicode.IsLetter(firstRune) && !unicode.IsLower(firstRune) {
		return ErrLowercaseStart
	}

	return nil
}
