package api

import (
	walletfacade "blockchain/app/wallet_server/api/wallet"
	"github.com/go-chi/chi/v5"
)

func (a *api) Router(r chi.Router) {
	r.Route("/api", func(r chi.Router) {
		walletfacade.RegisterRouter(r, a.cfg, a.walletUseCase)
	})
}
