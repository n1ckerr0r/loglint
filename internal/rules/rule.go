package rules

import "github.com/n1ckerr0r/loglint/internal/log_message"

type Rule interface {
	Name() string
	Check(message log_message.LogMessage) error
}
