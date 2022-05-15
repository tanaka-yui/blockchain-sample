package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func init() {
	log.SetPrefix("Block: ")
}

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash [32]byte) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
	}
}

func (b *Block) Print() {
	log.Println(fmt.Sprintf("timestamp:%d", b.timestamp))
	log.Println(fmt.Sprintf("nonce:%d", b.nonce))
	log.Println(fmt.Sprintf("previous_hash:%x", b.previousHash))
	log.Println(fmt.Sprintf("transactions:%s", b.transactions))
}

// 改ざんよう
func (b *Block) Mod() {
	b.nonce++
}

func (b *Block) Hash() [32]byte {
	m, err := json.Marshal(b)
	log.Println(fmt.Sprintf("marshal:%s", m))
	if err != nil {
		return [32]byte{}
	}
	return sha256.Sum256(m)
}

// MarshalJSON @Override
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64    `json:"timestamp"`
		Nonce        int      `json:"nonce"`
		PreviousHash [32]byte `json:"previous_hash"`
		Transactions []string `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		Transactions: b.transactions,
	})
}
