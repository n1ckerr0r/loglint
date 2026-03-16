package rules

import (
	"testing"

	"github.com/n1ckerr0r/loglint/internal/log_message"
)

func TestEnglishRule(t *testing.T) {

	rule := &EnglishRule{}

	msg := log_message.LogMessage{Text: "hello"}
	if err := rule.Check(msg); err != nil {
		t.Errorf("должно быть nil, получил %v", err)
	}

	msg = log_message.LogMessage{Text: "привет"}
	if err := rule.Check(msg); err != ErrNonEnglish {
		t.Errorf("должен быть ErrNonEnglish, получил %v", err)
	}
}

func TestLowercaseRule(t *testing.T) {

	rule := &LowercaseRule{}

	msg := log_message.LogMessage{Text: "hello"}
	if err := rule.Check(msg); err != nil {
		t.Errorf("должно быть nil, получил %v", err)
	}

	msg = log_message.LogMessage{Text: "Hello"}
	if err := rule.Check(msg); err != ErrLowercaseStart {
		t.Errorf("должен быть ErrLowercaseStart, получил %v", err)
	}

	msg = log_message.LogMessage{Text: "123 hello"}
	if err := rule.Check(msg); err != nil {
		t.Errorf("должно быть nil, получил %v", err)
	}
}

func TestSensitiveRule(t *testing.T) {

	rule := &SensitiveRule{
		keywords:      []string{"secret"},
		caseSensitive: true,
	}

	msg := log_message.LogMessage{Text: "this is secret"}
	if err := rule.Check(msg); err != ErrSensitiveData {
		t.Errorf("должен быть ErrSensitiveData, получил %v", err)
	}

	msg = log_message.LogMessage{Text: "this is SECRET"}
	if err := rule.Check(msg); err != nil {
		t.Errorf("должно быть nil, получил %v", err)
	}

	rule.caseSensitive = false

	msg = log_message.LogMessage{Text: "this is SECRET"}
	if err := rule.Check(msg); err != ErrSensitiveData {
		t.Errorf("должен быть ErrSensitiveData, получил %v", err)
	}
}

func TestSpecialCharsRule(t *testing.T) {

	rule := &SpecialCharsRule{}

	msg := log_message.LogMessage{Text: "hello world 123"}
	if err := rule.Check(msg); err != nil {
		t.Errorf("должно быть nil, получил %v", err)
	}

	msg = log_message.LogMessage{Text: "hello, world! how are you?"}
	if err := rule.Check(msg); err != nil {
		t.Errorf("должно быть nil, получил %v", err)
	}

	msg = log_message.LogMessage{Text: "hello § world"}
	if err := rule.Check(msg); err != ErrSpecialChars {
		t.Errorf("должен быть ErrSpecialChars, получил %v", err)
	}
}
