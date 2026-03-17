package rules

import "errors"

var (
	ErrLowercaseStart = errors.New("log message must start with logtest lowercase letter")
	ErrNonEnglish     = errors.New("log message must contain only english characters")
	ErrSpecialChars   = errors.New("log message contains forbidden characters")
	ErrSensitiveData  = errors.New("log message may contain sensitive data")
)
