package rules

import (
	"fmt"
	"go/token"
	"unicode"

	"github.com/n1ckerr0r/loglint/internal/log_message"
	"golang.org/x/tools/go/analysis"
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

func (specialCharRule *SpecialCharsRule) Check(message log_message.LogMessage) (error, *analysis.SuggestedFix) {

	for i, ch := range message.Text {

		if unicode.IsLetter(ch) ||
			unicode.IsDigit(ch) ||
			unicode.IsSpace(ch) {
			continue
		}

		pos := message.Pos + token.Pos(i+1)

		return ErrSpecialChars, &analysis.SuggestedFix{
			Message: fmt.Sprintf("remove forbidden character %q", ch),
			TextEdits: []analysis.TextEdit{
				{
					Pos: pos,
					End: pos + token.Pos(len(string(ch))),
					NewText: []byte(""),
				},
			},
		}
	}

	return nil, nil
}
