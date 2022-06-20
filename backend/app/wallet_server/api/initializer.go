package api

import (
	transactionrepository "blockchain/internal/repository/transaction"
	walletusecase "blockchain/internal/usecase/wallet"
	"blockchain/pkg/config"
)

func NewApi(
	cfg *config.Configuration,
) *api {
	a := &api{
		cfg: cfg,
	}
	initRepository(a)
	initUseCase(a)
	return a
}

func initRepository(a *api) {
	a.transactionRepository = transactionrepository.NewRepository()
}

func initUseCase(a *api) {
	a.walletUseCase = walletusecase.NewUseCase(a.cfg, a.transactionRepository)
}
