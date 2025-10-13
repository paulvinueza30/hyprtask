package logger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type CustomLogger struct {
	*slog.Logger
	tuiLog *slog.Logger
}

var Log *CustomLogger

func Init() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("Failed to create logs directory: " + err.Error())
	}

	// Main application log
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

	// TUI-specific log
	tuiLogFile := filepath.Join(logDir, "tui.log")
	tuiFile, err := os.OpenFile(tuiLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open TUI log file: " + err.Error())
	}

	tuiHandler := slog.NewTextHandler(tuiFile, &slog.HandlerOptions{
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

	Log = &CustomLogger{
		Logger: slog.New(handler),
		tuiLog: slog.New(tuiHandler),
	}
}

// Tui returns a logger that writes to tui.log
func (l *CustomLogger) Tui() *slog.Logger {
	return l.tuiLog
}
