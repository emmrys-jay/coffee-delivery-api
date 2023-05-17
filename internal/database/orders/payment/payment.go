package payment

type PaymentPayload struct {
	Email    string   `json:"email"`
	Amount   string   `json:"amount"`
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
}
