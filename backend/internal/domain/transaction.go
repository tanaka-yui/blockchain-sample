package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type BlockTransaction struct {
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

func NewBlockTransaction(sender string, recipient string, value float32) *BlockTransaction {
	return &BlockTransaction{sender, recipient, value}
}

func (t *BlockTransaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address      %s\n", t.SenderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address   %s\n", t.RecipientBlockchainAddress)
	fmt.Printf(" value                          %.1f\n", t.Value)
}

func (t *BlockTransaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.SenderBlockchainAddress,
		Recipient: t.RecipientBlockchainAddress,
		Value:     t.Value,
	})
}

func (t *BlockTransaction) UnmarshalJSON(data []byte) error {
	v := &struct {
		Sender    *string  `json:"sender_blockchain_address"`
		Recipient *string  `json:"recipient_blockchain_address"`
		Value     *float32 `json:"value"`
	}{
		Sender:    &t.SenderBlockchainAddress,
		Recipient: &t.RecipientBlockchainAddress,
		Value:     &t.Value,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
