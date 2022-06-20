package transactionusecase

import (
	"blockchain/internal/domain"
	transactionrepository "blockchain/internal/repository/transaction"
	"blockchain/pkg/config"
	"blockchain/pkg/customerrors"
	"blockchain/pkg/logger"
	"blockchain/pkg/utils/ecdsautil"
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
		ClearTransactionPool()
		GetChain() ([]*domain.Block, error)
		GetTransactionPool() []*domain.BlockTransaction
		CreateTransaction(input *CreatePutTransactionInput) error
		AddTransaction(input *CreatePutTransactionInput) bool
		GetAmount(input *GetAmountInput) (float32, error)
		Consensus() error
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

	uc.AddTransaction(&CreatePutTransactionInput{
		SenderBlockchainAddress: MiningSender,
		Value:                   MiningReward,
	})
	nonce := uc.proofOfWork()
	previousHash := uc.lastBlock().Hash()

	uc.CreateBlock(nonce, previousHash)

	for _, n := range uc.nodeIpAddresses {
		if err := uc.transactionRepository.Consensus(fmt.Sprintf("http://%s", n)); err != nil {
			logger.Logging.Error(fmt.Sprintf("Consensus error. %s", err.Error()))
		}
	}
}

func (uc *useCaseImpl) ClearTransactionPool() {
	uc.transactionPool = uc.transactionPool[:0]
}

func (uc *useCaseImpl) GetTransactionPool() []*domain.BlockTransaction {
	return uc.transactionPool
}

func (uc *useCaseImpl) CreateBlock(nonce int, previousHash [32]byte) {
	b := domain.NewBlock(nonce, previousHash, uc.transactionPool)
	uc.chain = append(uc.chain, b)

	uc.transactionPool = []*domain.BlockTransaction{}
	for _, n := range uc.nodeIpAddresses {
		if err := uc.transactionRepository.ClearTransaction(fmt.Sprintf("http://%s", n)); err != nil {
			logger.Logging.Warn(err.Error())
		}
	}
}

func (uc *useCaseImpl) GetChain() ([]*domain.Block, error) {
	return uc.chain, nil
}

func (uc *useCaseImpl) CreateTransaction(input *CreatePutTransactionInput) error {
	_ = uc.AddTransaction(input)

	for _, ip := range uc.nodeIpAddresses {
		err := uc.transactionRepository.PutTransaction(fmt.Sprintf("http://%s", ip), &transactionrepository.CreatePutRequest{
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

func (uc *useCaseImpl) AddTransaction(input *CreatePutTransactionInput) bool {
	t := domain.NewBlockTransaction(input.SenderBlockchainAddress, input.RecipientBlockchainAddress, input.Value)

	if input.SenderBlockchainAddress == MiningSender {
		uc.transactionPool = append(uc.transactionPool, t)
		return true
	}

	publicKey := ecdsautil.PublicKeyFromString(input.SenderPublicKey)
	signature := ecdsautil.SignatureFromString(input.Signature)

	if uc.verifyTransactionSignature(publicKey, signature, t) {
		//if uc.calculateTotalAmount(input.SenderBlockchainAddress) < input.Value {
		//	log.Println("ERROR: Not enough balance in a wallet")
		//	return false
		//}
		uc.transactionPool = append(uc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}

	return false
}

func (uc *useCaseImpl) GetAmount(input *GetAmountInput) (float32, error) {
	return uc.calculateTotalAmount(input.BlockchainAddress), nil
}

func (uc *useCaseImpl) Consensus() error {
	if uc.resolveConflicts() {
		return nil
	}
	return customerrors.NewError(customerrors.ErrCodeConflict, customerrors.WithMsg("Conflicts chain"))
}
