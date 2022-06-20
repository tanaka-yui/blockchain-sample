package walletfacade

type WalletOutput struct {
	PrivateKey        string `json:"privateKey"`
	PublicKey         string `json:"publicKey"`
	BlockchainAddress string `json:"blockchainAddress"`
}

type CreateTransactionInput struct {
	SenderPrivateKey           string `json:"senderPrivateKey"`
	SenderBlockchainAddress    string `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress string `json:"recipientBlockchainAddress"`
	SenderPublicKey            string `json:"senderPublicKey"`
	Value                      string `json:"value"`
}

type GetWalletResponse struct {
	Amount float32 `json:"amount"`
}
