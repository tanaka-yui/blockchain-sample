package api

import (
	"blockchain/pkg/config"
	"blockchain/pkg/logger"
	"blockchain/pkg/server/cron"
)

func NewApi(
	cfg *config.Configuration,
	cron *cron.Job,
) *api {
	a := &api{
		cfg:  cfg,
		cron: *cron,
	}
	initJob(a)
	return a
}

func initJob(a *api) {
	// "@every 1m" for debug
	a.cron.AddFunc("@every 10s", func() {
		logger.Logging.Info("cron job")
	})
}
