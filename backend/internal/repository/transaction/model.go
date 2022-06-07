package transactionrepository

type PutRequest struct {
	SenderBlockchainAddress    string  `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress string  `json:"recipientBlockchainAddress"`
	SenderPublicKey            string  `json:"senderPublicKey"`
	Value                      float32 `json:"value"`
	Signature                  string  `json:"signature"`
}
