package logtest

import "log/slog"

func goodLog() {
	slog.Info("server started")
}
