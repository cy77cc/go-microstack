// Package log log.go
package logx

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	gklog "github.com/go-kit/log"
)

// Level 定义日志级别
type Level int

// Option 定义Logger选项
type Option func(*Logger)

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志记录器结构体
type Logger struct {
	level        Level
	logger       gklog.Logger
	mu           sync.Mutex
	timeFormat   string
	host         string
	pid          int
	colorEnabled bool
}

// 全局默认logger实例
var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(INFO, os.Stdout, WithCaller())
}

// WithCaller 返回一个启用caller功能的Option
func WithCaller() Option {
	return func(l *Logger) {
		l.logger = gklog.With(l.logger, "caller", gklog.Caller(4))
	}
}

// NewLogger 创建新的Logger实例
func NewLogger(level Level, output io.Writer, opts ...Option) *Logger {
	// 使用go-kit/log创建基础logger
	gkLogger := gklog.NewLogfmtLogger(output)

	// 添加时间戳
	gkLogger = gklog.With(gkLogger, "ts", gklog.TimestampFormat(
		func() time.Time { return time.Now() },
		"2006-01-02 15:04:05",
	))

	// 应用选项
	logger := &Logger{
		level:        level,
		logger:       gkLogger,
		timeFormat:   "2006-01-02 15:04:05",
		pid:          os.Getpid(),
		colorEnabled: false,
	}
	if hn, err := os.Hostname(); err == nil {
		logger.host = hn
	}
	if logger.host != "" {
		logger.logger = gklog.With(logger.logger, "host", logger.host)
	}
	logger.logger = gklog.With(logger.logger, "pid", logger.pid)
	// 通过环境变量启用颜色
	if v := os.Getenv("LOG_COLOR"); v == "1" || v == "true" || v == "TRUE" {
		logger.colorEnabled = true
	}

	for _, opt := range opts {
		opt(logger)
	}

	return logger
}

// SetLevel 设置日志级别
func SetLevel(level Level) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.level = level
}

// SetOutput 设置输出目标
func SetOutput(output io.Writer) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()

	// 重新创建go-kit logger
	gkLogger := gklog.NewLogfmtLogger(output)
	gkLogger = gklog.With(gkLogger, "ts", gklog.TimestampFormat(
		func() time.Time { return time.Now() },
		defaultLogger.timeFormat,
	))
	defaultLogger.logger = gkLogger
}

// SetTimeFormat 设置时间格式
func SetTimeFormat(format string) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	defaultLogger.timeFormat = format

	// 更新时间戳格式
	gkLogger := gklog.NewLogfmtLogger(gklog.NewSyncWriter(os.Stdout))
	gkLogger = gklog.With(gkLogger, "ts", gklog.TimestampFormat(
		func() time.Time { return time.Now() },
		format,
	))
	defaultLogger.logger = gkLogger
}

// Debug 记录debug级别日志
func Debug(v ...interface{}) {
	if defaultLogger.level <= DEBUG {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(DEBUG), "msg", fmt.Sprint(v...))
	}
}

// Debugf 格式化记录debug级别日志
func Debugf(format string, v ...interface{}) {
	if defaultLogger.level <= DEBUG {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(DEBUG), "msg", fmt.Sprintf(format, v...))
	}
}

// Info 记录info级别日志
func Info(v ...interface{}) {
	if defaultLogger.level <= INFO {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(INFO), "msg", fmt.Sprint(v...))
	}
}

// Infof 格式化记录info级别日志
func Infof(format string, v ...interface{}) {
	if defaultLogger.level <= INFO {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(INFO), "msg", fmt.Sprintf(format, v...))
	}
}

// Warn 记录warn级别日志
func Warn(v ...interface{}) {
	if defaultLogger.level <= WARN {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(WARN), "msg", fmt.Sprint(v...))
	}
}

// Warnf 格式化记录warn级别日志
func Warnf(format string, v ...interface{}) {
	if defaultLogger.level <= WARN {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(WARN), "msg", fmt.Sprintf(format, v...))
	}
}

// Error 记录error级别日志
func Error(v ...interface{}) {
	if defaultLogger.level <= ERROR {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(ERROR), "msg", fmt.Sprint(v...))
	}
}

// Errorf 格式化记录error级别日志
func Errorf(format string, v ...interface{}) {
	if defaultLogger.level <= ERROR {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(ERROR), "msg", fmt.Sprintf(format, v...))
	}
}

// Fatal 记录fatal级别日志并退出程序
func Fatal(v ...interface{}) {
	if defaultLogger.level <= FATAL {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(FATAL), "msg", fmt.Sprint(v...))
		os.Exit(1)
	}
}

// Fatalf 格式化记录fatal级别日志并退出程序
func Fatalf(format string, v ...interface{}) {
	if defaultLogger.level <= FATAL {
		defaultLogger.logger.Log("level", defaultLogger.coloredLevel(FATAL), "msg", fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}

// WithFields 创建带字段的日志记录器
func WithFields(fields map[string]interface{}) *FieldLogger {
	// 将fields转换为go-kit/log使用的键值对
	keyvals := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		keyvals = append(keyvals, k, v)
	}

	return &FieldLogger{
		logger: gklog.With(defaultLogger.logger, keyvals...),
		fields: fields,
	}
}

// FieldLogger 带字段的日志记录器
type FieldLogger struct {
	logger gklog.Logger
	fields map[string]interface{}
}

// Debug 记录带字段的debug级别日志
func (fl *FieldLogger) Debug(v ...interface{}) {
	if defaultLogger.level <= DEBUG {
		fl.logger.Log("level", defaultLogger.coloredLevel(DEBUG), "msg", fmt.Sprint(v...))
	}
}

// Info 记录带字段的info级别日志
func (fl *FieldLogger) Info(v ...interface{}) {
	if defaultLogger.level <= INFO {
		fl.logger.Log("level", defaultLogger.coloredLevel(INFO), "msg", fmt.Sprint(v...))
	}
}

// Warn 记录带字段的warn级别日志
func (fl *FieldLogger) Warn(v ...interface{}) {
	if defaultLogger.level <= WARN {
		fl.logger.Log("level", defaultLogger.coloredLevel(WARN), "msg", fmt.Sprint(v...))
	}
}

// Error 记录带字段的error级别日志
func (fl *FieldLogger) Error(v ...interface{}) {
	if defaultLogger.level <= ERROR {
		fl.logger.Log("level", defaultLogger.coloredLevel(ERROR), "msg", fmt.Sprint(v...))
	}
}

func (l *Logger) coloredLevel(level Level) string {
	s := level.String()
	if !l.colorEnabled {
		return s
	}
	const (
		reset   = "\033[0m"
		blue    = "\033[34m"
		green   = "\033[32m"
		yellow  = "\033[33m"
		red     = "\033[31m"
		redBold = "\033[1;31m"
	)
	switch level {
	case DEBUG:
		return blue + s + reset
	case INFO:
		return green + s + reset
	case WARN:
		return yellow + s + reset
	case ERROR:
		return red + s + reset
	case FATAL:
		return redBold + s + reset
	default:
		return s
	}
}
