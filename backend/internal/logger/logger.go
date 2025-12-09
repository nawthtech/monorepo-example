package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

type DefaultLogger struct {
	logger *slog.Logger
}

var (
	Stdout     *slog.Logger
	Stderr     *slog.Logger
	globalLogger Logger
)

func Init(env string) {
	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	if env == "development" {
		Stdout = slog.New(slog.NewTextHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewTextHandler(os.Stderr, opts))
	} else {
		opts.AddSource = true
		Stdout = slog.New(slog.NewJSONHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	globalLogger = &DefaultLogger{logger: Stdout}
}

func (l *DefaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}
func (l *DefaultLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}
func (l *DefaultLogger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}
func (l *DefaultLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
func (l *DefaultLogger) With(args ...any) Logger {
	return &DefaultLogger{logger: l.logger.With(args...)}
}

func Debug(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Debug(ctx, msg, args...)
}
func Info(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Info(ctx, msg, args...)
}
func Warn(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Warn(ctx, msg, args...)
}
func Error(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development")
	}
	globalLogger.Error(ctx, msg, args...)
}
func With(args ...any) Logger {
	if globalLogger == nil {
		Init("development")
	}
	return globalLogger.With(args...)
}

// helpers...
func ErrAttr(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "")
	}
	return slog.String("error", err.Error())
}
func DurationAttr(d time.Duration) slog.Attr {
	return slog.Duration("duration", d)
}
func TimestampAttr() slog.Attr {
	return slog.String("timestamp", time.Now().Format(time.RFC3339))
}
func RequestAttr(method, path string, statusCode int, duration time.Duration) slog.Attr {
	return slog.Group("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
	)
}
func MemoryUsageAttr() slog.Attr {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return slog.Group("memory",
		slog.Any("alloc", m.Alloc),
		slog.Any("sys", m.Sys),
		slog.Uint64("num_gc", uint64(m.NumGC)),
	)
}