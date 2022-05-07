package logger

import (
	"blockchain/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"sync"
)

type (
	logging interface {
		Debug(msg string, fields ...zap.Field)
		Info(msg string, fields ...zap.Field)
		Warn(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
		Fatal(msg string, fields ...zap.Field)

		Sync()
		GetZapLogger() *zap.Logger
	}
	loggingImpl struct {
		zapLogger *zap.Logger
		logLevel  zapcore.Level
	}
)

var (
	// Workaround for sync error
	// https://github.com/uber-go/zap/issues/880#issuecomment-731261906
	ignoreSyncError bool
	Logging         logging
	once            sync.Once
)

func (l *loggingImpl) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *loggingImpl) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *loggingImpl) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *loggingImpl) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *loggingImpl) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *loggingImpl) Sync() {
	if err := l.zapLogger.Sync(); err != nil {
		if ignoreSyncError {
			return
		}
		log.Printf("failed to sync log: %v", err)
	}
}

func (l *loggingImpl) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

func LoadLogger(c *config.Configuration) {
	once.Do(func() {
		logLevel := unmarshalLogLevel(c.Logger.LogLevel)
		zapLogger := newZapLogger(c.Logger.LogEncoding, logLevel)
		Logging = &loggingImpl{
			zapLogger: zapLogger,
			logLevel:  logLevel,
		}
		if err := zapLogger.Sync(); err != nil {
			ignoreSyncError = true
		}
		log.Printf("load logger, env: %s, logLevel: %v", c.System.Env, logLevel)
	})
}

func newZapLogger(encoding string, logLevel zapcore.Level) *zap.Logger {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         encoding,
		EncoderConfig:    getEncodeConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("init logging configs error: %v", err)
	}
	return logger
}

func getEncodeConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// see: https://github.com/uber-go/zap/blob/425214515ff452748375576b20c82524849177c6/zapcore/level.go#L126-L146
func unmarshalLogLevel(text string) zapcore.Level {
	switch text {
	case "debug", "DEBUG":
		return zap.DebugLevel
	case "info", "INFO", "": // make the zero value useful
		return zap.InfoLevel
	case "warn", "WARN":
		return zap.WarnLevel
	case "error", "ERROR":
		return zap.ErrorLevel
	case "dpanic", "DPANIC":
		return zap.DPanicLevel
	case "panic", "PANIC":
		return zap.PanicLevel
	case "fatal", "FATAL":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}
