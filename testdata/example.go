package testdata

import (
	"log/slog"
)

func main() {
	slog.Info("это сообщение на русском")
	slog.Info("It's wrong message")
	slog.Info("Oh my gosh, this is a special char #@#@##")
	slog.Info("this password so sensitive(((")
}
