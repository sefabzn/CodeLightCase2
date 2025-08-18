package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/db"
	"app/internal/handlers"
	"app/internal/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading .env file: %v", err)
		log.Printf("Will use system environment variables instead")
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Load configuration
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Supabase client
	database := db.NewSupabaseClient(config.SupabaseURL, config.SupabaseAnonKey, config.SupabaseServiceKey)

	// Create Echo instance
	e := echo.New()

	// Setup all routes and middleware
	handlers.SetupRoutes(e, database)

	// Setup graceful shutdown
	go func() {
		if err := e.Start(config.GetAddr()); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Server startup failed: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Gracefully shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	database.Close()
	log.Println("Server exited")
}
