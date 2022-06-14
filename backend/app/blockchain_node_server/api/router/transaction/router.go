package transaction

import (
	transactionusecase "blockchain/internal/usecase/transaction"
	"blockchain/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterRouter(
	router chi.Router,
	cfg *config.Configuration,
	transactionUseCase transactionusecase.UseCase,
) {
	f := &facade{
		validator:          validator.New(),
		cfg:                cfg,
		transactionUseCase: transactionUseCase,
	}
	router.Route("/blockchain", func(a chi.Router) {
		a.Route("/chains", func(b chi.Router) {
			b.Get("/", f.GetChain)
		})
		a.Route("/transactions", func(b chi.Router) {
			b.Get("/", f.GetTransactionPool)
			b.Post("/", f.CreateTransaction)
			b.Put("/", f.PutTransaction)
		})
	})
}
