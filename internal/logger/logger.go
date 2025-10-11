package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var Log *slog.Logger

func Init() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Failed to create logs directory: " + err.Error())
	}

	logFile := filepath.Join(logDir, "hyprtask.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if idx := strings.Index(a.Value.String(), "hyprtask/"); idx != -1 {
					a.Value = slog.StringValue(a.Value.String()[idx:])
				}
			}
			return a
		},
	})
	Log = slog.New(handler)
}
