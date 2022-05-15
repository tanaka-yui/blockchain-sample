package chain

import (
	"fmt"
	"log"
	"strings"

	"blockchain/internal/domain/block"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

type Blockchain struct {
	transactionPool []string
	chain           []*block.Block
}

func NewBlockchain() *Blockchain {
	b := &block.Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *block.Block {
	b := block.NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) Print() {
	for i, b := range bc.chain {
		log.Println(fmt.Sprintf("%s Chain %d %s", strings.Repeat("=", 25), i, strings.Repeat("=", 25)))
		b.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bc *Blockchain) LastBlock() *block.Block {
	return bc.chain[len(bc.chain)-1]
}
