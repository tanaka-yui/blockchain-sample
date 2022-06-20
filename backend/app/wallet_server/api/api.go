package api

import (
	transactionrepository "blockchain/internal/repository/transaction"
	walletusecase "blockchain/internal/usecase/wallet"
	"blockchain/pkg/config"
)

type (
	api struct {
		cfg     *config.Configuration
		gateway string

		transactionRepository transactionrepository.Repository

		walletUseCase walletusecase.UseCase
	}
)
