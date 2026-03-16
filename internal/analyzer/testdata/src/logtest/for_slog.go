package logtest

import "log/slog"

func slogExample() {
	slog.Info("Bad Message") // want "lowercase"
}
