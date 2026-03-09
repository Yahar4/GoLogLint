package sensitive

import "log/slog"

func main() {
	slog.Info("user logged in")  // ok
	slog.Info("password: 12345") // want "log-message cant contain any sensitive data"
}
