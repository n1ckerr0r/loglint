package panicrule

import "log/slog"

func panicExample() {
	slog.Info("message") // want "loglint internal error"
}
