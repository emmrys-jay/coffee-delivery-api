package handlers

import (
	"net/http"
	"strconv"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database/models"
	"github.com/emmrys-jay/coffee-delivery-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	UserService *services.UserService
	Validator   *validator.Validate
}

func NewUserHandler(userService *services.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Validator:   validator,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.CreateUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.Validator.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	retUser, err := h.UserService.CreateUser(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Status: true, Message: "User created successfully", Data: retUser})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	user, err := h.UserService.GetUserByID(c, uint(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "User retrieved successfully", Data: user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	var user models.UserUpdate
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.Validator.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.UserService.UpdateUser(c, uint(idInt), &user); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "User updated successfully", Data: user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.UserService.DeleteUser(c, uint(idInt)); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusNoContent, models.Response{Status: true, Message: "User deleted successfully", Data: nil})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.UserService.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Users retrieved successfully", Data: users})
}

func (h *UserHandler) Login(c *gin.Context) {
	var credentials models.LoginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	if err := h.Validator.Struct(credentials); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	token, err := h.UserService.Login(c, credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{Status: false, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Status: true, Message: "Login successful", Data: gin.H{"token": token}})
}
