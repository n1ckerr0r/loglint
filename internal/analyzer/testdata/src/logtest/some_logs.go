package logtest

import "log/slog"

func someLogs() {
	slog.Info("bad message")
	slog.Info("Bad message") // want "lowercase"
}
