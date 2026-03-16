package lowercase

import "log/slog"

func slogExample() {
	slog.Info("Bad Message") // want "lowercase"
}
