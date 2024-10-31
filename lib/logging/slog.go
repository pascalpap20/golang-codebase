package logging

import (
	"log"
	"log/slog"
	"os"

	"example.com/example/config"
)

var slogJson *slog.Logger = slog.Default()

func InitSlogLogger(c *config.Config) {
	var opts *slog.HandlerOptions
	if c.IsDevelopment() || c.IsTesting() {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
	}
	slogJson = slog.New(slog.NewJSONHandler(os.Stdout, opts))
}

func Debug(msg string, keysAndValues ...interface{}) {
	slogJson.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	slogJson.Info(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	slogJson.Warn(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	slogJson.Error(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	slogJson.Error(msg, keysAndValues...)
	log.Fatal(msg)
}
