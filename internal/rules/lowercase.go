package rules

import (
	"go/token"
	"unicode"
	"unicode/utf8"

	"github.com/n1ckerr0r/loglint/internal/log_message"
	"golang.org/x/tools/go/analysis"
)

// Это правило считает лог, начинающийся с заглавной буквы некорректным
type LowercaseRule struct{}

func NewLowercaseRule() *LowercaseRule {
	return &LowercaseRule{}
}

func (lowercaseRule *LowercaseRule) Name() string {
	return "LowercaseFirstRune"
}

func (lowercaseRule *LowercaseRule) Check(message log_message.LogMessage) (error, *analysis.SuggestedFix) {
	if message.Text == "" {
		return nil, nil
	}

	firstRune, size := utf8.DecodeRuneInString(message.Text)

	if unicode.IsLetter(firstRune) && !unicode.IsLower(firstRune) {

		fixed := string(unicode.ToLower(firstRune)) + message.Text[size:]

		fix := &analysis.SuggestedFix{
			Message: "make message lowercase",
			TextEdits: []analysis.TextEdit{
				{
					Pos: message.Pos,
					End: message.Pos + token.Pos(len(message.Text)+2),
					NewText: []byte(`"` + fixed + `"`),
				},
			},
		}

		return ErrLowercaseStart, fix
	}

	return nil, nil
}
