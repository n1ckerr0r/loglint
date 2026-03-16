package rules

import (
	"unicode"

	"github.com/n1ckerr0r/loglint/internal/log_message"
)

// Специальными символами считаются все символы, которые не являются буквами, цифрами или пробелом
// и не входят в список разрещенных символов (пользователю можно расширить список)
type SpecialCharsRule struct{}

func NewSpecialCharsRule() *SpecialCharsRule {
	return &SpecialCharsRule{}
}

func (specialCharRule *SpecialCharsRule) Name() string {
	return "SpecialChars"
}

func (specialCharRule *SpecialCharsRule) Check(logMessage log_message.LogMessage) error {

	allowedSpecial := map[rune]bool{
		'.': true, ',': true, '!': true, '?': true,
		':': true, ';': true, '-': true, '\'': true,
		'"': true, '(': true, ')': true, '[': true, ']': true,
		'{': true, '}': true, '@': true, '#': true,
		'$': true, '%': true, '^': true, '&': true,
		'*': true, '+': true, '=': true, '<': true,
		'>': true, '/': true, '\\': true, '|': true,
		'~': true, '`': true,
	}

	for _, ch := range logMessage.Text {

		if unicode.IsLetter(ch) ||
			unicode.IsDigit(ch) ||
			unicode.IsSpace(ch) {
			continue
		}

		if allowedSpecial[ch] {
			continue
		}

		return ErrSpecialChars
	}

	return nil
}
