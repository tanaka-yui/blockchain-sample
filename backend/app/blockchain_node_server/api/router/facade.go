package router

import (
	"blockchain/pkg/config"
	"blockchain/pkg/network/response"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type facade struct {
	validator *validator.Validate
	cfg       *config.Configuration
}

func (f *facade) GetTransactionPool(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response.OKWithMsg(ctx, w, &struct {
		Ok bool `json:"ok"`
	}{
		Ok: true,
	})
}
