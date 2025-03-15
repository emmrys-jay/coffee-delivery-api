package util

import (
	"errors"

	"github.com/shopspring/decimal"
)

func ParseDecimal(amount string) (decimal.Decimal, error) {
	num, err := decimal.NewFromString(amount)
	if err != nil {
		return decimal.Decimal{}, errors.New("Invalid amount passed in")
	}

	return num, nil
}
