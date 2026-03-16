package logtest

import "log/slog"

func extractorFail() {
	msg := "Bad Message"
	slog.Info(msg) // want "lowercase"
}
