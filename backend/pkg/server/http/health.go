package http

import (
	"blockchain/pkg/network/response"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const HealthCheckPath = "/api/health"

type healthCheckResponse struct {
	Message string
}

func registerHealth(router chi.Router) {
	router.Get(HealthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		response.OKWithMsg(
			r.Context(),
			w,
			healthCheckResponse{
				Message: "OK",
			},
		)
	})
}
