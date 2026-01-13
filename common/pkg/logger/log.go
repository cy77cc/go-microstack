package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	level  zap.AtomicLevel
)

func Init(levelStr string) error {
	level = zap.NewAtomicLevel()
	levelStr = strings.ToLower(levelStr)
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return err
	}

	cfg := zap.Config{
		Level:       level,
		Development: true,
		Encoding:    "json", // 生产推荐 json
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "ts",
			LevelKey:    "level",
			MessageKey:  "msg",
			CallerKey:   "caller",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	logger, err = cfg.Build()
	return err
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(format string, a ...any) {
	logger.Info(fmt.Sprintf(format, a...))
}

func Level() string {
	return level.String()
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(format string, a ...any) {
	logger.Warn(fmt.Sprintf(format, a...))
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(format string, a ...any) {
	logger.Error(fmt.Sprintf(format, a...))
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(format string, a ...any) {
	logger.Debug(fmt.Sprintf(format, a...))
}