package router

import (
	"blockchain/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func RegisterRouter(
	router chi.Router,
	cfg *config.Configuration,
) {
	f := &facade{
		validator: validator.New(),
		cfg:       cfg,
	}
	router.Route("/blockchain", func(a chi.Router) {
		a.Get("/pool/transactions", f.GetTransactionPool)
	})
}
