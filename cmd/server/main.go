package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emmrys-jay/coffee-delivery-api/internal/database"
	"github.com/emmrys-jay/coffee-delivery-api/internal/database/repository"
	"github.com/emmrys-jay/coffee-delivery-api/internal/handlers"
	"github.com/emmrys-jay/coffee-delivery-api/internal/routes"
	"github.com/emmrys-jay/coffee-delivery-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load configurations: %v", err)
	}

	// Connect to the database
	db, err := database.ConnectAndAutoMigrate()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Initialize the repository
	coffeeRepo := repository.NewCoffeeRepository(db)
	coffeeService := services.NewCoffeeService(coffeeRepo)
	coffeeHandler := handlers.NewCoffeeHandler(coffeeService)

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Set up the Gin router
	router := gin.Default()
	routes.SetupRoutes(router, coffeeHandler, userHandler)

	// Start the server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Staring server on " + srv.Addr)

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
