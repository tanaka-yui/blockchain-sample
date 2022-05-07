package api

import (
	"blockchain/pkg/config"
	"blockchain/pkg/server/cron"
)

type (
	api struct {
		cfg  *config.Configuration
		cron cron.Job
	}
)
