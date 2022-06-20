package transactionusecase

type Transaction struct {
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

type CreatePutTransactionInput struct {
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	SenderPublicKey            string
	Value                      float32
	Signature                  string
}

type GetAmountInput struct {
	BlockchainAddress string
}
