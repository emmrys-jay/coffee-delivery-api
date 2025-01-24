package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CoffeeHandler struct {
	service   *services.CoffeeService
	validator *validator.Validate
}

func NewCoffeeHandler(service *services.CoffeeService) *CoffeeHandler {
	return &CoffeeHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *CoffeeHandler) CreateCoffee(c *gin.Context) {
	var coffee models.Coffee
	if err := c.ShouldBindJSON(&coffee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(coffee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coffee.CreatedAt = time.Now()
	coffee.UpdatedAt = time.Now()

	if err := h.service.CreateCoffee(c, &coffee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, coffee)
}

func (h *CoffeeHandler) GetCoffee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	coffee, err := h.service.GetCoffeeByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coffee not found"})
		return
	}

	c.JSON(http.StatusOK, coffee)
}

func (h *CoffeeHandler) UpdateCoffee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var coffee models.Coffee
	if err := c.ShouldBindJSON(&coffee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(coffee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coffee.Id = id
	coffee.UpdatedAt = time.Now()

	if err := h.service.UpdateCoffee(c, &coffee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coffee)
}

func (h *CoffeeHandler) DeleteCoffee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteCoffee(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *CoffeeHandler) ListCoffees(c *gin.Context) {
	coffees, err := h.service.ListCoffees(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coffees)
}
