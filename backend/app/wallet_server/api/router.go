package api

import "github.com/go-chi/chi/v5"

func (a *api) Router(r chi.Router) {
	r.Route("/api", func(r chi.Router) {

	})
}
