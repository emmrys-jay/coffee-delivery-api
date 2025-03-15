package handlers

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/services"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// OrderHandler represents the HTTP handler for order-related requests
type OrderHandler struct {
	service  *services.OrderService
	validate *validator.Validate
}

// NewOrderHandler creates a new OrderHandler instance
func NewOrderHandler(svc *services.OrderService, vld *validator.Validate) *OrderHandler {
	return &OrderHandler{
		svc,
		vld,
	}
}

// CreateOrder handles the creation of a new order
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
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

	order, err := h.service.PlaceOrder(c, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Status: true, Message: "Order created successfully", Data: order})
}

// GetOrder handles fetching a single order by ID
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "Invalid order ID", Data: nil})
		return
	}

	order, err := h.service.GetOrder(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Order fetched successfully", Data: order})
}

// ListUsersOrders handles fetching a single order by ID
func (h *OrderHandler) ListUsersOrders(c *gin.Context) {
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

	order, err := h.service.ListUserOrders(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Order fetched successfully", Data: order})
}

// UpdateOrder handles updating an existing order
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "Invalid order ID", Data: nil})
		return
	}

	var req models.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if req.Status == "" {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "No status was specified", Data: nil})
		return
	}

	order, err := h.service.UpdateOrderStatus(c, uint(id), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Order updated successfully", Data: order})
}

// CancelOrder handles canceling an order by ID
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "Invalid order ID", Data: nil})
		return
	}

	order, err := h.service.CancelOrder(c, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Order canceled successfully", Data: order})
}
