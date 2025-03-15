package platform

type PaymentProvider string

const (
	PAYSTACK PaymentProvider = "Paystack"
)

func (pp PaymentProvider) String() string {
	return string(pp)
}

type InitTransactionRequest struct {
	Amount    string `json:"amount,omitempty"`
	Email     string `json:"email,omitempty"`
	Reference string `json:"reference,omitempty"`
}

type InitTransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}
