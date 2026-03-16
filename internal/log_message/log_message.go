package log_message

import (
	"go/token"
)

// Представляет лог-запись, извлечённую из AST кода
// Добавлен для расширяемости пока что используется слабо
type LogMessage struct {
	Text       string
	Pos        token.Pos
	LoggerType string
	Level      string
}

func NewLogMessage(text string, pos token.Pos, loggerType string, level string) *LogMessage {
	return &LogMessage{
		Text:       text,
		Pos:        pos,
		LoggerType: loggerType,
		Level:      level,
	}
}
