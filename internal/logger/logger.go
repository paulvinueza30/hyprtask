package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

func Init() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
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