package walletusecase

type GetAmountInput struct {
	BlockChainAddress string
}

type CreateTransactionInput struct {
	SenderPrivateKey           string `validate:"required" json:"sender_private_key"`
	SenderBlockchainAddress    string `validate:"required" json:"sender_blockchain_address"`
	RecipientBlockchainAddress string `validate:"required" json:"recipient_blockchain_address"`
	SenderPublicKey            string `validate:"required" json:"sender_public_key"`
	Value                      string `validate:"required" json:"value"`
}
