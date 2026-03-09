package specialchars

import "log/slog"

func main() {
	slog.Info("normal message") // ok
	slog.Info("bad!!!")         // want "log-message cant contain any special symbols"
}
