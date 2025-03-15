package util

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func amountValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	decimal, err := decimal.NewFromString(value)
	if err != nil {
		return false
	}

	if decimal.IsNegative() || decimal.IsZero() {
		return false
	}
	return true
}

func NewValidator() *validator.Validate {
	validate := validator.New()

	// Registering the new rule "sig" for significant
	err := validate.RegisterValidation("sig", amountValidator)
	if err != nil {
		fmt.Println("Error registering custom validation :", err.Error())
	}

	return validate
}
