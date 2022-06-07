package api

import (
	"blockchain/app/blockchain_node_server/api/router/transaction"
	"github.com/go-chi/chi/v5"
)

func (a *api) Router(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		transaction.RegisterRouter(r, a.cfg, a.transactionUseCase)
	})
}
