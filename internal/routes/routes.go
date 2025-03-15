package routes

import (
	"github.com/emmrys-jay/coffee-delivery-api/internal/handlers"
	"github.com/emmrys-jay/coffee-delivery-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	coffeeHandler *handlers.CoffeeHandler,
	userHandler *handlers.UserHandler,
	orderHandler *handlers.OrderHandler,
	trxHandler *handlers.TransactionHandler,
) {
	// Public routes
	router.POST("/login", userHandler.Login)
	router.POST("/users", userHandler.CreateUser)

	// Authenticated user routes
	auth := router.Group("/")
	auth.Use(middlewares.UserAuthMiddleware())
	{
		auth.GET("/coffees", coffeeHandler.ListCoffees)
		auth.GET("/coffees/:id", coffeeHandler.GetCoffee)

		auth.POST("/orders", orderHandler.CreateOrder)
		auth.GET("/orders/:id", orderHandler.GetOrder)
		auth.GET("/orders", orderHandler.ListUsersOrders)
		auth.PATCH("/orders/cancel", orderHandler.CancelOrder)
		auth.POST("/orders/pay", trxHandler.InitiatePayment)
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

		admin.PATCH("/orders/:id", orderHandler.UpdateOrder)
	}

	// Webhook routes
	webhook := router.Group("/")
	webhook.Use(middlewares.IPWhitelistMiddleware())
	{
		webhook.POST("/webhook/paystack", trxHandler.HandlePaystackWebhook)
	}
}
