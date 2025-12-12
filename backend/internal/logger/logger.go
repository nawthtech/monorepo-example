package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Logger الواجهة الرئيسية
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Fatal(ctx context.Context, msg string, args ...any)
	
	// Health & Monitoring specific methods
	Health(ctx context.Context, component string, status string, duration time.Duration, args ...any)
	Metric(ctx context.Context, name string, value float64, args ...any)
	Alert(ctx context.Context, level, source, message string, args ...any)
	
	With(args ...any) Logger
	WithGroup(name string) Logger
}

// DefaultLogger التطبيق الأساسي
type DefaultLogger struct {
	logger *slog.Logger
}

var (
	Stdout *slog.Logger
	Stderr *slog.Logger

	globalLogger Logger
	appName      string = "nawthtech"
	appVersion   string = "1.0.0"
)

// Init تهيئة النظام
func Init(env, name, version string) {
	if name != "" {
		appName = name
	}
	if version != "" {
		appVersion = version
	}

	level := slog.LevelInfo
	if env == "development" || env == "debug" {
		level = slog.LevelDebug
	} else if env == "test" {
		level = slog.LevelWarn
	}

	opts := &slog.HandlerOptions{
		Level:       level,
		AddSource:   env == "development" || env == "test",
		ReplaceAttr: replaceAttr,
	}

	// اختيار شكل السجل حسب البيئة
	if env == "development" {
		Stdout = slog.New(slog.NewTextHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewTextHandler(os.Stderr, opts))
	} else {
		Stdout = slog.New(slog.NewJSONHandler(os.Stdout, opts))
		Stderr = slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	// تهيئة السجل العالمي مع معلومات التطبيق
	globalLogger = &DefaultLogger{
		logger: Stdout.With(
			slog.String("app", appName),
			slog.String("version", appVersion),
			slog.String("env", env),
		),
	}

	globalLogger.Info(context.Background(), "Logger initialized", 
		slog.String("level", level.String()),
		slog.String("format", getFormat(env)),
	)
}

// replaceAttr دالة لتعديل السمات
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{
			Key:   "timestamp",
			Value: slog.StringValue(time.Now().Format(time.RFC3339Nano)),
		}
	}
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		if source != nil {
			return slog.Attr{
				Key:   "source",
				Value: slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line)),
			}
		}
	}
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		return slog.Attr{
			Key:   "level",
			Value: slog.StringValue(level.String()),
		}
	}
	return a
}

// getFormat تحديد شكل السجل
func getFormat(env string) string {
	if env == "development" {
		return "text"
	}
	return "json"
}

// ================================
// التطبيقات الأساسية
// ================================

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

func (l *DefaultLogger) Fatal(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

// ================================
// دوال Health & Monitoring
// ================================

func (l *DefaultLogger) Health(ctx context.Context, component string, status string, duration time.Duration, args ...any) {
	attrs := []any{
		slog.String("component", component),
		slog.String("status", status),
		slog.Duration("check_duration", duration),
		slog.Time("checked_at", time.Now()),
	}
	attrs = append(attrs, args...)
	
	l.logger.InfoContext(ctx, "health_check", attrs...)
}

func (l *DefaultLogger) Metric(ctx context.Context, name string, value float64, args ...any) {
	attrs := []any{
		slog.String("metric_name", name),
		slog.Float64("metric_value", value),
		slog.Time("measured_at", time.Now()),
	}
	attrs = append(attrs, args...)
	
	l.logger.InfoContext(ctx, "metric_collected", attrs...)
}

func (l *DefaultLogger) Alert(ctx context.Context, level, source, message string, args ...any) {
	attrs := []any{
		slog.String("alert_level", level),
		slog.String("alert_source", source),
		slog.String("alert_message", message),
		slog.Time("alert_time", time.Now()),
	}
	attrs = append(attrs, args...)
	
	switch level {
	case "critical", "fatal":
		l.logger.ErrorContext(ctx, "system_alert", attrs...)
	case "warning", "error":
		l.logger.WarnContext(ctx, "system_alert", attrs...)
	default:
		l.logger.InfoContext(ctx, "system_alert", attrs...)
	}
}

func (l *DefaultLogger) With(args ...any) Logger {
	return &DefaultLogger{logger: l.logger.With(args...)}
}

func (l *DefaultLogger) WithGroup(name string) Logger {
	return &DefaultLogger{logger: l.logger.WithGroup(name)}
}

// ================================
// دوال عامة
// ================================

func Debug(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Debug(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Info(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Warn(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Error(ctx, msg, args...)
}

func Fatal(ctx context.Context, msg string, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Fatal(ctx, msg, args...)
}

func Health(ctx context.Context, component string, status string, duration time.Duration, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Health(ctx, component, status, duration, args...)
}

func Metric(ctx context.Context, name string, value float64, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Metric(ctx, name, value, args...)
}

func Alert(ctx context.Context, level, source, message string, args ...any) {
	if globalLogger == nil {
		Init("development", "", "")
	}
	globalLogger.Alert(ctx, level, source, message, args...)
}

// ================================
// دوال مساعدة
// ================================

func ErrAttr(err error) slog.Attr {
	if err == nil {
		return slog.String("error", "")
	}
	return slog.String("error", err.Error())
}

func DurationAttr(duration time.Duration) slog.Attr {
	return slog.Duration("duration", duration)
}

func TimestampAttr() slog.Attr {
	return slog.String("timestamp", time.Now().Format(time.RFC3339Nano))
}

func ComponentAttr(component string) slog.Attr {
	return slog.String("component", component)
}

func StatusAttr(status string) slog.Attr {
	return slog.String("status", status)
}

func RequestAttr(method, path string, statusCode int, duration time.Duration) slog.Attr {
	return slog.Group("request",
		slog.String("method", method),
		slog.String("path", path),
		slog.Int("status_code", statusCode),
		slog.Duration("duration", duration),
		slog.String("user_agent", getUserAgent()),
		slog.String("ip_address", getIPAddress()),
	)
}

func HealthAttr(component, status string, duration time.Duration) slog.Attr {
	return slog.Group("health",
		slog.String("component", component),
		slog.String("status", status),
		slog.Duration("check_duration", duration),
		slog.Time("checked_at", time.Now()),
	)
}

func MetricAttr(name string, value float64, unit string) slog.Attr {
	return slog.Group("metric",
		slog.String("name", name),
		slog.Float64("value", value),
		slog.String("unit", unit),
		slog.Time("measured_at", time.Now()),
	)
}

func SystemAttr() slog.Attr {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	return slog.Group("system",
		slog.String("go_version", runtime.Version()),
		slog.String("os", runtime.GOOS),
		slog.String("arch", runtime.GOARCH),
		slog.Int("num_cpu", runtime.NumCPU()),
		slog.Int("goroutines", runtime.NumGoroutine()),
		slog.String("memory_alloc", formatMemory(memStats.Alloc)),
		slog.String("memory_sys", formatMemory(memStats.Sys)),
		slog.String("memory_heap", formatMemory(memStats.HeapAlloc)),
	)
}

// ================================
// دوال مساعدة خاصة
// ================================

func formatMemory(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func getIPAddress() string {
	// TODO: يمكنك إضافة منطق لاستخراج IP
	return "unknown"
}

func getUserAgent() string {
	// TODO: يمكنك إضافة منطق لاستخراج User-Agent
	return "unknown"
}

// ================================
// دوال الإدارة
// ================================

func GetGlobalLogger() Logger {
	if globalLogger == nil {
		Init("development", "", "")
	}
	return globalLogger
}

func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

func GetAppName() string {
	return appName
}

func GetAppVersion() string {
	return appVersion
}

// ================================
// دوال لبناء السجلات المنظمة
// ================================

func NewRequestLogger(method, path string) Logger {
	if globalLogger == nil {
		Init("development", "", "")
	}
	return globalLogger.With(
		slog.String("method", method),
		slog.String("path", path),
		slog.String("request_id", generateRequestID()),
	)
}

func NewHealthLogger(component string) Logger {
	if globalLogger == nil {
		Init("development", "", "")
	}
	return globalLogger.With(
		slog.String("component", component),
		slog.String("monitoring", "true"),
	)
}

func generateRequestID() string {
	return fmt.Sprintf("req_%d_%d", time.Now().UnixNano(), os.Getpid())
}

// ================================
// تهيئة الافتراضية (للتطوير)
// ================================

func init() {
	// تهيئة تلقائية للتطوير
	if os.Getenv("APP_ENV") == "" {
		Init("development", "nawthtech", "1.0.0")
	}
}