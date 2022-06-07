package transaction

type PutRequest struct {
	SenderBlockchainAddress    string  `json:"senderBlockchainAddress,omitempty"`
	RecipientBlockchainAddress string  `json:"recipientBlockchainAddress,omitempty"`
	SenderPublicKey            string  `json:"senderPublicKey,omitempty"`
	Value                      float32 `json:"value,omitempty"`
	Signature                  string  `json:"signature,omitempty"`
}
