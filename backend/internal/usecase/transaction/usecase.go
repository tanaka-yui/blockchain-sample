package transactionusecase

import (
	"blockchain/internal/domain"
	transactionrepository "blockchain/internal/repository/transaction"
	"blockchain/pkg/config"
	"blockchain/pkg/logger"
	"blockchain/pkg/utils/neighbor"
	"fmt"
	"log"
	"strings"
	"sync"
)

type (
	UseCase interface {
		FineNeighbors()
		DebugTransactionPool()
		Mining()
		GetChain() ([]*domain.Block, error)
		CreateTransaction(input *CreatePutTransactionInput) error
		AddTransaction(input *CreatePutTransactionInput) error
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
	myPort uint16,
	cfg *config.Configuration,
	transactionRepository transactionrepository.Repository,
) UseCase {
	return &useCaseImpl{
		myPort:                myPort,
		cfg:                   cfg,
		transactionRepository: transactionRepository,
		nodeIpAddresses:       []string{},
		transactionPool:       []*domain.BlockTransaction{},
		chain:                 []*domain.Block{},
	}
}

const (
	BlockchainPortRangeStart = 5001
	BlockchainPortRangeEnd   = 5003
	NeighborIpRangeStart     = 0
	NeighborIpRangeEnd       = 1

	MiningSender = "THE BLOCKCHAIN"
	MiningReward = 1.0
)

func (uc *useCaseImpl) setNodeIpAddresses(nodeIpAddresses []string) {
	uc.nodeIpAddresses = nodeIpAddresses
}

func (uc *useCaseImpl) FineNeighbors() {
	uc.lock.Lock()
	defer uc.lock.Unlock()

	uc.nodeIpAddresses = neighbor.FindNeighbors(
		neighbor.GetHost(), uc.myPort,
		NeighborIpRangeStart, NeighborIpRangeEnd,
		BlockchainPortRangeStart, BlockchainPortRangeEnd)
	log.Printf("%v", uc.nodeIpAddresses)
}

func (uc *useCaseImpl) DebugTransactionPool() {
	logger.Logging.Info(strings.Repeat("-", 40))
	for _, v := range uc.chain {
		v.Print()
	}
	logger.Logging.Info(strings.Repeat("-", 40))
}

func (uc *useCaseImpl) Mining() {
	uc.lockMining.Lock()
	defer uc.lockMining.Unlock()

	if err := uc.AddTransaction(&CreatePutTransactionInput{
		SenderBlockchainAddress: MiningSender,
		Value:                   MiningReward,
	}); err != nil {
		logger.Logging.Warn(err.Error())
	}

	nonce := uc.proofOfWork()
	previousHash := uc.lastBlock().Hash()

	uc.CreateBlock(nonce, previousHash)
}

func (uc *useCaseImpl) CreateBlock(nonce int, previousHash [32]byte) {
	b := domain.NewBlock(nonce, previousHash, uc.transactionPool)
	uc.chain = append(uc.chain, b)
}

func (uc *useCaseImpl) GetChain() ([]*domain.Block, error) {
	return uc.chain, nil
}

func (uc *useCaseImpl) CreateTransaction(input *CreatePutTransactionInput) error {
	_ = uc.AddTransaction(input)

	for _, ip := range uc.nodeIpAddresses {
		err := uc.transactionRepository.PutTransaction(fmt.Sprintf("http://%s/api/blockchain/transactions", ip), &transactionrepository.PutRequest{
			SenderBlockchainAddress:    input.SenderBlockchainAddress,
			RecipientBlockchainAddress: input.RecipientBlockchainAddress,
			SenderPublicKey:            input.SenderPublicKey,
			Value:                      input.Value,
			Signature:                  input.Signature,
		})
		if err != nil {
			logger.Logging.Warn(err.Error())
		}
	}
	return nil
}

func (uc *useCaseImpl) AddTransaction(input *CreatePutTransactionInput) error {
	t := domain.NewBlockTransaction(input.SenderBlockchainAddress, input.RecipientBlockchainAddress, input.Value)
	uc.transactionPool = append(uc.transactionPool, t)

	return nil
}
