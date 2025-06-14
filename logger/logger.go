package logger

import (
	"log/slog"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func Initiate(level slog.Level) {
	// Setup lumberjack for log rotation
	rotatingFile := &lumberjack.Logger{
		Filename:   filepath.Join("log", "probe.log"), // Log file path
		MaxSize:    10,                                // Max megabytes before rotation
		MaxBackups: 10,                                // Max old log files to keep
		MaxAge:     20,                                // Max days to retain a log file
		Compress:   true,                              // Compress old logs
	}

	// Create a slog handler writing to lumberjack logger
	handler := slog.NewTextHandler(rotatingFile, &slog.HandlerOptions{
		Level: level, //slog.LevelDebug, // Set log level here
	})

	slog.SetDefault(slog.New(handler))
}
