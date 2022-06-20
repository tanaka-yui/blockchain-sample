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

	response.OKWithMsg(ctx, w, &GetChainResponse{
		Chain: result,
	})
}

func (f *facade) GetTransactionPool(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := f.transactionUseCase.GetTransactionPool()

	response.OKWithMsg(ctx, w, &GetTransactionResponse{
		Transactions: result,
	})
}

func (f *facade) ClearTransactionPool(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f.transactionUseCase.ClearTransactionPool()

	response.OK(ctx, w)
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

	if success := f.transactionUseCase.AddTransaction(&input); !success {
		response.BadRequest(ctx, w)
		return
	}

	response.OK(ctx, w)
}

func (f *facade) GetAmount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := f.transactionUseCase.GetAmount(&transactionusecase.GetAmountInput{
		BlockchainAddress: json.Query(r, "blockchainAddress"),
	})
	if err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OKWithMsg(ctx, w, &GetWalletResponse{
		Amount: result,
	})
}

func (f *facade) Consensus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := f.transactionUseCase.Consensus(); err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OK(ctx, w)
}
