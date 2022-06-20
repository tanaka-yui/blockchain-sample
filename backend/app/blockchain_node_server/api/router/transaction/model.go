package transaction

import "blockchain/internal/domain"

type PutRequest struct {
	SenderBlockchainAddress    string  `json:"senderBlockchainAddress,omitempty"`
	RecipientBlockchainAddress string  `json:"recipientBlockchainAddress,omitempty"`
	SenderPublicKey            string  `json:"senderPublicKey,omitempty"`
	Value                      float32 `json:"value,omitempty"`
	Signature                  string  `json:"signature,omitempty"`
}

type GetChainResponse struct {
	Chain []*domain.Block `json:"chain"`
}

type GetWalletResponse struct {
	Amount float32 `json:"amount"`
}

type GetTransactionResponse struct {
	Transactions []*domain.BlockTransaction `json:"transactions"`
}
