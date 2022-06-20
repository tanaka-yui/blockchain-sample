package walletfacade

import (
	walletusecase "blockchain/internal/usecase/wallet"
	"blockchain/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterRouter(
	router chi.Router,
	cfg *config.Configuration,
	walletUseCase walletusecase.UseCase,
) {
	f := &facade{
		validator:     validator.New(),
		cfg:           cfg,
		walletUseCase: walletUseCase,
	}
	router.Route("/blockchain", func(a chi.Router) {
		a.Route("/wallet", func(b chi.Router) {
			b.Post("/", f.CreateWallet)
			b.Get("/", f.GetWallet)
		})
		a.Route("/transactions", func(b chi.Router) {
			b.Post("/", f.CreateTransaction)
		})
	})
}
