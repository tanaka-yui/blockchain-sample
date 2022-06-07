package api

import (
	transactionrepository "blockchain/internal/repository/transaction"
	transactionusecase "blockchain/internal/usecase/transaction"
	"blockchain/pkg/config"
	"blockchain/pkg/server/cron"
)

func NewApi(
	port uint16,
	cfg *config.Configuration,
	cron *cron.Job,
) *api {
	a := &api{
		port: port,
		cfg:  cfg,
		cron: *cron,
	}
	initRepository(a)
	initUseCase(a)
	initJob(a)
	return a
}

func initRepository(a *api) {
	a.transactionRepository = transactionrepository.NewRepository()
}

func initUseCase(a *api) {
	a.transactionUseCase = transactionusecase.NewUseCase(a.port, a.cfg, a.transactionRepository)
}

func initJob(a *api) {
	// "@every 1m" for debug
	a.cron.AddFunc("@every 10s", a.transactionUseCase.FineNeighbors)
	a.cron.AddFunc("@every 10s", a.transactionUseCase.DebugTransactionPool)
}
