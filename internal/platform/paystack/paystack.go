package paystack

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/emmrys-jay/coffee-delivery-api/internal/platform"
	"github.com/emmrys-jay/coffee-delivery-api/util"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type PaystackService struct {
	SecretKey string
	BaseUrl   string
}

func NewPaystackService() *PaystackService {
	return &PaystackService{
		SecretKey: os.Getenv("PAYSTACK_SECRET"),
		BaseUrl:   os.Getenv("PAYSTACK_BASEURL"),
	}
}

func (ps *PaystackService) InitiateTransaction(email string, amount decimal.Decimal, reference string) (platform.InitTransactionResponse, error) {
	if ps.BaseUrl == "" || ps.SecretKey == "" {
		return platform.InitTransactionResponse{}, errors.New("paystack is not configured")
	}

	url := ps.BaseUrl + "/transaction/initialize"
	headers := map[string]string{
		"Authorization": "Bearer " + ps.SecretKey,
	}
	body := platform.InitTransactionRequest{
		Amount:    convertToSubunit(amount).String(),
		Email:     email,
		Reference: reference,
	}

	respBody, err := util.MakePOSTRequest(url, headers, body)
	if err != nil {
		logrus.Error("Error InitializeTransaction: ", err)
		return platform.InitTransactionResponse{}, err
	}

	var response platform.InitTransactionResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		logrus.Error("Error InitializeTransaction, unmarshaling result: ", err)
		return platform.InitTransactionResponse{}, err
	}

	return response, nil
}

func convertToSubunit(amount decimal.Decimal) decimal.Decimal {
	return amount.Mul(decimal.New(100, 0)) // multiply by 100
}
