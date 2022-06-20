package walletusecase

import (
	"blockchain/internal/domain"
	transactionrepository "blockchain/internal/repository/transaction"
	"blockchain/pkg/config"
	"blockchain/pkg/customerrors"
	"blockchain/pkg/utils/ecdsautil"
	"strconv"
	"sync"
)

type (
	UseCase interface {
		CreateWallet() (*domain.Wallet, error)
		GetAmount(input *GetAmountInput) (float32, error)
		CreateTransaction(input *CreateTransactionInput) error
	}
	useCaseImpl struct {
		myPort                uint16
		cfg                   *config.Configuration
		transactionRepository transactionrepository.Repository

		lock            sync.Mutex
		lockMining      sync.Mutex
		nodeIpAddresses []string
		transactionPool []*domain.BlockTransaction
		chain           []*domain.Block
	}
)

func NewUseCase(
	cfg *config.Configuration,
	transactionRepository transactionrepository.Repository,
) UseCase {
	return &useCaseImpl{
		cfg:                   cfg,
		transactionRepository: transactionRepository,
		nodeIpAddresses:       []string{},
		transactionPool:       []*domain.BlockTransaction{},
		chain:                 []*domain.Block{},
	}
}

func (uc *useCaseImpl) CreateWallet() (*domain.Wallet, error) {
	return domain.NewWallet(), nil
}

func (uc *useCaseImpl) GetAmount(input *GetAmountInput) (float32, error) {
	res, err := uc.transactionRepository.GetAmount(uc.cfg.System.Http.NodeGateway, &transactionrepository.GetAmountInput{
		BlockchainAddress: input.BlockChainAddress,
	})
	if err != nil {
		return 0.0, customerrors.NewError(customerrors.ErrUnexpected, customerrors.WithError(err))
	}

	return res, nil
}

func (uc *useCaseImpl) CreateTransaction(input *CreateTransactionInput) error {
	publicKey := ecdsautil.PublicKeyFromString(input.SenderPublicKey)
	privateKey := ecdsautil.PrivateKeyFromString(input.SenderPrivateKey, publicKey)
	value, err := strconv.ParseFloat(input.Value, 32)
	if err != nil {
		return customerrors.NewError(customerrors.ErrUnexpected, customerrors.WithError(err))
	}
	value32 := float32(value)
	transaction := domain.NewWalletTransaction(privateKey, publicKey, input.SenderBlockchainAddress, input.RecipientBlockchainAddress, value32)
	signature := transaction.GenerateSignature()
	signatureStr := signature.String()
	if err := uc.transactionRepository.CreateTransaction(uc.cfg.System.Http.NodeGateway, &transactionrepository.CreatePutRequest{
		SenderBlockchainAddress:    input.SenderBlockchainAddress,
		RecipientBlockchainAddress: input.RecipientBlockchainAddress,
		SenderPublicKey:            input.SenderPublicKey,
		Value:                      value32,
		Signature:                  signatureStr,
	}); err != nil {
		return customerrors.NewError(customerrors.ErrUnexpected, customerrors.WithError(err))
	}
	return nil
}
