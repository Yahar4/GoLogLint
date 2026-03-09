package lowercase

import "log/slog"

func main() {
	slog.Info("correct") // ok
	slog.Info("Wrong")   // want "log-message must be in lowercase"
}
