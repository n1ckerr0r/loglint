package rules

import (
	"github.com/n1ckerr0r/loglint/internal/log_message"
	"golang.org/x/tools/go/analysis"
)

type Rule interface {
	Name() string
	Check(message log_message.LogMessage) (error, *analysis.SuggestedFix)
}
