package platform

import "github.com/shopspring/decimal"

type Provider interface {
	InitiateTransaction(email string, amount decimal.Decimal, reference string) (InitTransactionResponse, error)
}
