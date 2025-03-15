package handlers

import (
	"net/http"
	"strconv"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CoffeeHandler struct {
	service   *services.CoffeeService
	validator *validator.Validate
}

func NewCoffeeHandler(service *services.CoffeeService, validator *validator.Validate) *CoffeeHandler {
	return &CoffeeHandler{
		service:   service,
		validator: validator,
	}
}

func (h *CoffeeHandler) CreateCoffee(c *gin.Context) {
	var coffee models.CreateCoffee
	if err := c.ShouldBindJSON(&coffee); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.validator.Struct(coffee); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	retCoffee, err := h.service.CreateCoffee(c, &coffee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Status: true, Message: "Coffee created successfully", Data: retCoffee})
}

func (h *CoffeeHandler) GetCoffee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "Invalid ID", Data: nil})
		return
	}

	coffee, err := h.service.GetCoffeeByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Status: false, Message: "Coffee not found", Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Coffee retrieved successfully", Data: coffee})
}

func (h *CoffeeHandler) UpdateCoffee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "Invalid ID", Data: nil})
		return
	}

	var coffee models.UpdateCoffee
	if err := c.ShouldBindJSON(&coffee); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.validator.Struct(coffee); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.service.UpdateCoffee(c, uint(id), &coffee); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Coffee updated successfully", Data: coffee})
}

func (h *CoffeeHandler) DeleteCoffee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: "Invalid ID", Data: nil})
		return
	}

	if err := h.service.DeleteCoffee(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusNoContent, models.Response{Status: true, Message: "Coffee deleted successfully", Data: nil})
}

func (h *CoffeeHandler) ListCoffees(c *gin.Context) {
	coffees, err := h.service.ListCoffees(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Coffees retrieved successfully", Data: coffees})
}
