package rules

import (
	"unicode"

	"github.com/n1ckerr0r/loglint/internal/log_message"
	"golang.org/x/tools/go/analysis"
)

// Правило запрещает все естественные языки кроме английского
type EnglishRule struct{}

func NewEnglishRule() *EnglishRule {
	return &EnglishRule{}
}

func (englishRule *EnglishRule) Name() string {
	return "EnglishLanguage"
}

func (englishRule *EnglishRule) Check(logMessage log_message.LogMessage) (error, *analysis.SuggestedFix) {
	for _, char := range logMessage.Text {

		if char <= unicode.MaxASCII {
			continue
		}

		// Если символ - буква, но не ASCII, значит это другой язык
		if unicode.IsLetter(char) {
			return ErrNonEnglish, nil
		}

		continue
	}
	return nil, nil
}
