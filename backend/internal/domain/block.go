package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	log.SetPrefix("Block: ")
}

type Block struct {
	nonce        int
	previousHash [32]byte
	timestamp    int64
	transactions []*BlockTransaction
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*BlockTransaction) *Block {
	return &Block{
		nonce:        nonce,
		previousHash: previousHash,
		timestamp:    time.Now().UnixNano(),
		transactions: transactions,
	}
}

func (b *Block) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	log.Println(fmt.Sprintf("hash:%x", b.Hash()))
	log.Println(fmt.Sprintf("timestamp:%d", b.timestamp))
	log.Println(fmt.Sprintf("nonce:%d", b.nonce))
	log.Println(fmt.Sprintf("previous_hash:%x", b.previousHash))
	fmt.Printf("%s\n", strings.Repeat("-", 40))
}

func (b *Block) Hash() [32]byte {
	m, err := json.Marshal(b)
	//log.Println(fmt.Sprintf("marshal:%s", m))
	if err != nil {
		return [32]byte{}
	}
	return sha256.Sum256(m)
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64               `json:"timestamp"`
		Nonce        int                 `json:"nonce"`
		PreviousHash string              `json:"previous_hash"`
		Transactions []*BlockTransaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PreviousHash: fmt.Sprintf("%x", b.previousHash),
		Transactions: b.transactions,
	})
}

func (b *Block) UnmarshalJSON(data []byte) error {
	var previousHash string
	v := &struct {
		Timestamp    *int64               `json:"timestamp"`
		Nonce        *int                 `json:"nonce"`
		PreviousHash *string              `json:"previous_hash"`
		Transactions *[]*BlockTransaction `json:"transactions"`
	}{
		Timestamp:    &b.timestamp,
		Nonce:        &b.nonce,
		PreviousHash: &previousHash,
		Transactions: &b.transactions,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	ph, _ := hex.DecodeString(*v.PreviousHash)
	copy(b.previousHash[:], ph[:32])
	return nil
}
