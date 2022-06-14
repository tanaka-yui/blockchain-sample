package transaction

import (
	transactionusecase "blockchain/internal/usecase/transaction"
	"blockchain/pkg/config"
	"blockchain/pkg/network/request/json"
	"blockchain/pkg/network/response"
	"context"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type facade struct {
	validator          *validator.Validate
	cfg                *config.Configuration
	transactionUseCase transactionusecase.UseCase
}

func (f *facade) GetChain(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := f.transactionUseCase.GetChain()
	if err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OKWithMsg(ctx, w, result)
}

func (f *facade) GetTransactionPool(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response.OKWithMsg(ctx, w, &struct {
		Ok bool `json:"ok"`
	}{
		Ok: true,
	})
}

func (f *facade) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var body PutRequest
	if err := json.Bind(&body, r); err != nil {
		response.BadRequest(ctx, w, response.WithError(err))
		return
	}

	input := transactionusecase.CreatePutTransactionInput{
		SenderBlockchainAddress:    body.SenderBlockchainAddress,
		RecipientBlockchainAddress: body.RecipientBlockchainAddress,
		SenderPublicKey:            body.SenderPublicKey,
		Value:                      body.Value,
		Signature:                  body.Signature,
	}

	if err := f.transactionUseCase.CreateTransaction(&input); err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OK(ctx, w)
}

func (f *facade) PutTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var body PutRequest
	if err := json.Bind(&body, r); err != nil {
		response.BadRequest(ctx, w, response.WithError(err))
		return
	}

	input := transactionusecase.CreatePutTransactionInput{
		SenderBlockchainAddress:    body.SenderBlockchainAddress,
		RecipientBlockchainAddress: body.RecipientBlockchainAddress,
		SenderPublicKey:            body.SenderPublicKey,
		Value:                      body.Value,
		Signature:                  body.Signature,
	}

	if err := f.transactionUseCase.AddTransaction(&input); err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OK(ctx, w)
}
