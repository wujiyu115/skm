package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger *slog.Logger

func init() {
	defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, nil))
}

type Config struct {
	Dir        string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Debug      bool
}

func Setup(cfg Config) error {
	if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
		return err
	}

	maxSize := cfg.MaxSizeMB
	if maxSize <= 0 {
		maxSize = 10
	}
	maxBackups := cfg.MaxBackups
	if maxBackups <= 0 {
		maxBackups = 5
	}
	maxAge := cfg.MaxAgeDays
	if maxAge <= 0 {
		maxAge = 30
	}

	lj := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Dir, "skm.log"),
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LocalTime:  true,
	}

	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}

	w := io.MultiWriter(os.Stderr, lj)
	defaultLogger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(defaultLogger)

	return nil
}

func Info(msg string, args ...any)  { defaultLogger.Info(msg, args...) }
func Warn(msg string, args ...any)  { defaultLogger.Warn(msg, args...) }
func Error(msg string, args ...any) { defaultLogger.Error(msg, args...) }
func Debug(msg string, args ...any) { defaultLogger.Debug(msg, args...) }
