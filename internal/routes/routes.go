package routes

import (
	"github.com/emmrys-jay/coffee-delivery-api/internal/handlers"
	"github.com/emmrys-jay/coffee-delivery-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, coffeeHandler *handlers.CoffeeHandler, userHandler *handlers.UserHandler) {
	// Public routes
	router.POST("/login", userHandler.Login)
	router.POST("/users", userHandler.CreateUser)

	// Authenticated user routes
	auth := router.Group("/")
	auth.Use(middlewares.UserAuthMiddleware())
	{
		auth.GET("/coffees", coffeeHandler.ListCoffees)
		auth.GET("/coffees/:id", coffeeHandler.GetCoffee)
	}

	// Admin routes
	admin := router.Group("/")
	admin.Use(middlewares.AdminAuthMiddleware())
	{
		admin.POST("/coffees", coffeeHandler.CreateCoffee)
		admin.PUT("/coffees/:id", coffeeHandler.UpdateCoffee)
		admin.DELETE("/coffees/:id", coffeeHandler.DeleteCoffee)

		admin.GET("/users", userHandler.ListUsers)
		admin.GET("/users/:id", userHandler.GetUser)
		admin.PUT("/users/:id", userHandler.UpdateUser)
		admin.DELETE("/users/:id", userHandler.DeleteUser)
	}
}
