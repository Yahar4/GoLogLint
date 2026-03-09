package russian

import "log/slog"

func main() {
	slog.Info("english") // ok
	slog.Info("привет")  // want "log-message must be in english"
}
