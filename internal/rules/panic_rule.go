package rules

import "github.com/n1ckerr0r/loglint/internal/log_message"

// Правило создано для тестов, на случай создания кастомных правил, нужно обработать случае, если эти правило вызывают панику
type PanicRule struct{}

func (PanicRule) Name() string {
	return "panic-rule"
}

func (PanicRule) Check(log_message.LogMessage) error {
	panic("rule panic")
}
