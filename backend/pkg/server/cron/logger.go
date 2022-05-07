package cron

import "go.uber.org/zap"

type jobLogger struct {
	logger *zap.Logger
}

func (j *jobLogger) Info(msg string, keysAndValues ...interface{}) {
	j.logger.Info(msg, zap.Any("attributes", keysAndValues))
}

func (j *jobLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	j.logger.Error(msg, zap.Any("attributes", keysAndValues), zap.Error(err))
}
