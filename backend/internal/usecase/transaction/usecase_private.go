package transactionusecase

import (
	"blockchain/internal/domain"
	"blockchain/pkg/logger"
	"blockchain/pkg/utils/ecdsautil"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

const (
	MiningDifficulty = 3
)

func (uc *useCaseImpl) createBlock(nonce int, previousHash [32]byte) *domain.Block {
	b := domain.NewBlock(nonce, previousHash, uc.transactionPool)
	uc.chain = append(uc.chain, b)
	uc.transactionPool = []*domain.BlockTransaction{}
	return b
}

func (uc *useCaseImpl) lastBlock() *domain.Block {
	if len(uc.chain) == 0 {
		return nil
	}
	return uc.chain[len(uc.chain)-1]
}

func (uc *useCaseImpl) copyTransactionPool() []*domain.BlockTransaction {
	transactions := make([]*domain.BlockTransaction, 0)
	for _, t := range uc.transactionPool {
		transactions = append(transactions,
			domain.NewBlockTransaction(t.SenderBlockchainAddress,
				t.RecipientBlockchainAddress,
				t.Value))
	}
	return transactions
}

func (uc *useCaseImpl) verifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, s *ecdsautil.Signature, t *domain.BlockTransaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256(m)
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (uc *useCaseImpl) validProof(nonce int, previousHash [32]byte, transactions []*domain.BlockTransaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := domain.NewBlock(nonce, previousHash, transactions)
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (uc *useCaseImpl) validChain(chain []*domain.Block) bool {
	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.PreviousHash() != preBlock.Hash() {
			return false
		}

		if !uc.validProof(b.Nonce(), b.PreviousHash(), b.Transactions(), MiningDifficulty) {
			return false
		}

		preBlock = b
		currentIndex += 1
	}
	return true
}

func (uc *useCaseImpl) proofOfWork() int {
	transactions := uc.copyTransactionPool()
	previousHash := uc.lastBlock().Hash()
	nonce := 0
	for !uc.validProof(nonce, previousHash, transactions, MiningDifficulty) {
		nonce += 1
	}
	return nonce
}

func (uc *useCaseImpl) debugChainPrint() {
	for i, block := range uc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (uc *useCaseImpl) calculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range uc.chain {
		for _, t := range b.Transactions() {
			value := t.Value
			if blockchainAddress == t.RecipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.SenderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (uc *useCaseImpl) resolveConflicts() bool {
	var longestChain []*domain.Block = nil
	maxLength := len(uc.chain)

	for _, n := range uc.nodeIpAddresses {
		res, err := uc.transactionRepository.GetChain(fmt.Sprintf("http://%s", n))
		if err != nil {
			logger.Logging.Error("Error getting chain from node", zap.Error(err))
			continue
		}
		if len(res) > maxLength && uc.validChain(res) {
			maxLength = len(res)
			longestChain = res
		}
	}

	if longestChain != nil {
		uc.chain = longestChain
		logger.Logging.Info("Resolve conflicts replaced")
		return true
	}
	logger.Logging.Info("Resolve conflicts not replaced")
	return false
}
