package test

import (
	"blockchain/internal/domain/block"
	"blockchain/internal/domain/chain"
	"blockchain/internal/domain/wallet"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
)

func Test_Block(t *testing.T) {
	blockChain := chain.NewBlockchain()
	blockChain.Print()
	blockChain.CreateBlock(5, blockChain.LastBlock().Hash())
	blockChain.Print()
	blockChain.CreateBlock(2, blockChain.LastBlock().Hash())
	blockChain.Print()
}

func Test_Wallet(t *testing.T) {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())

	b := block.NewBlock(1, [32]byte{})

	//const TestMessage = "Just some test message..."
	//hash := sha256.Sum256([]byte(TestMessage))

	hash := b.Hash()
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey(), hash[:])
	if err != nil {
		log.Println("Error signing message: ", err)
		log.Fatal(err)
	}
	// 改ざんする
	//b.Mod()

	// ハッシュ値が改ざんされてないかチェックする
	hash = b.Hash()
	ok := ecdsa.Verify(w.PublicKey(), hash[:], r, s)
	log.Printf("verify: %v", ok)
}
