package walletfacade

import (
	walletusecase "blockchain/internal/usecase/wallet"
	"blockchain/pkg/config"
	"blockchain/pkg/network/request/json"
	"blockchain/pkg/network/response"
	"context"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type facade struct {
	validator     *validator.Validate
	cfg           *config.Configuration
	walletUseCase walletusecase.UseCase
}

func (f *facade) CreateWallet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	result, err := f.walletUseCase.CreateWallet()
	if err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OKWithMsg(ctx, w, &WalletOutput{
		BlockchainAddress: result.BlockchainAddress(),
		PublicKey:         result.PublicKeyStr(),
		PrivateKey:        result.PrivateKeyStr(),
	})
}

func (f *facade) GetWallet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	result, err := f.walletUseCase.GetAmount(&walletusecase.GetAmountInput{
		BlockChainAddress: json.Query(r, "blockchainAddress"),
	})
	if err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OKWithMsg(ctx, w, &GetWalletResponse{
		Amount: result,
	})
}

func (f *facade) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var body CreateTransactionInput
	if err := json.Bind(&body, r); err != nil {
		response.BadRequest(ctx, w, response.WithError(err))
		return
	}

	input := walletusecase.CreateTransactionInput{
		SenderBlockchainAddress:    body.SenderBlockchainAddress,
		RecipientBlockchainAddress: body.RecipientBlockchainAddress,
		SenderPublicKey:            body.SenderPublicKey,
		SenderPrivateKey:           body.SenderPrivateKey,
		Value:                      body.Value,
	}

	if err := f.walletUseCase.CreateTransaction(&input); err != nil {
		response.ErrorResponse(ctx, w, err)
		return
	}

	response.OK(ctx, w)
}
