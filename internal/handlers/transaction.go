package handlers

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TransactionHandler struct {
	service  *services.TransactionService
	validate *validator.Validate
}

func NewTransactionHandler(service *services.TransactionService, validator *validator.Validate) *TransactionHandler {
	return &TransactionHandler{
		service:  service,
		validate: validator,
	}
}

// CreateOrder handles the creation of a new order
func (h *TransactionHandler) InitiatePayment(c *gin.Context) {
	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	claims, _ := c.Get("claims")
	userId, ok := claims.(jwt.MapClaims)["user_id"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "could not get user id from request", Data: nil})
		return
	}

	id, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "could not get user id from request", Data: nil})
		return
	}

	trx, err := h.service.Initiate(c, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Transaction created successfully", Data: trx})
}

func (h *TransactionHandler) HandlePaystackWebhook(c *gin.Context) {
	// Verify and process Paystack webhook
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
