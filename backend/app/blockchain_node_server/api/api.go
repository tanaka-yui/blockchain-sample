package api

import (
	transactionrepository "blockchain/internal/repository/transaction"
	transactionusecase "blockchain/internal/usecase/transaction"
	"blockchain/pkg/config"
	"blockchain/pkg/server/cron"
)

type (
	api struct {
		port uint16
		cfg  *config.Configuration
		cron cron.Job

		transactionRepository transactionrepository.Repository

		transactionUseCase transactionusecase.UseCase
	}
)
