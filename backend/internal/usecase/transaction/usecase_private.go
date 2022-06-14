package transactionusecase

import (
	"blockchain/internal/domain"
	"fmt"
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

func (uc *useCaseImpl) validProof(nonce int, previousHash [32]byte, transactions []*domain.BlockTransaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := domain.NewBlock(nonce, previousHash, transactions)
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
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
