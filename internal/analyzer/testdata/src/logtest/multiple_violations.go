package logtest

import "log/slog"

func multipleViolations() {
	slog.Info("Password leaked!!!") // want "sensitive"
}
