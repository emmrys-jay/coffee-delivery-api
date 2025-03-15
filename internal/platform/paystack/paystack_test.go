package paystack

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestInitiateTransaction(t *testing.T) {
	ps := &PaystackService{
		BaseUrl:   "http://localhost:8080",
		SecretKey: "test_secret_key",
	}

	// Mock server to simulate Paystack API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/transaction/initialize", r.URL.Path)
		require.Equal(t, "Bearer test_secret_key", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status": true, "message": "Authorization URL created", "data": {"authorization_url": "https://paystack.com/pay/abc123"}}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	ps.BaseUrl = server.URL

	email := "test@example.com"
	amount := decimal.NewFromFloat(100.00)
	reference := "test_reference"

	response, err := ps.InitiateTransaction(email, amount, reference)
	require.NoError(t, err)
	require.Equal(t, "https://paystack.com/pay/abc123", response.Data.AuthorizationURL)
}

func TestInitiateTransaction_Integration(t *testing.T) {
	ps := &PaystackService{
		BaseUrl:   "https://api.paystack.co",
		SecretKey: "sk_test_ae76d5df2f8464aac34f1013c0fd02059c5d1a49", // Test Secret Key from paystack
	}

	email := "test@example.com"
	amount := decimal.NewFromFloat(100.00)
	reference := uuid.New().String()

	response, err := ps.InitiateTransaction(email, amount, reference)
	require.NoError(t, err)

	logrus.Info(response)
}
