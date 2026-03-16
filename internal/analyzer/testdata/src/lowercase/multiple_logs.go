package lowercase

import "log/slog"

func multipleLogs() {
	slog.Info("bad message")
	slog.Info("Bad message") // want "lowercase"
}
