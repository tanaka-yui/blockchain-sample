package transactionrepository

import "blockchain/internal/domain"

type CreatePutRequest struct {
	SenderBlockchainAddress    string  `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress string  `json:"recipientBlockchainAddress"`
	SenderPublicKey            string  `json:"senderPublicKey"`
	Value                      float32 `json:"value"`
	Signature                  string  `json:"signature"`
}

type GetAmountInput struct {
	BlockchainAddress string `json:"blockchainAddress"`
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}

type GetChainResponse struct {
	Chain []*domain.Block `json:"chain"`
}
