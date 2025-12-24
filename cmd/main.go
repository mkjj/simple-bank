package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"simple_bank/api/handler"
	"simple_bank/config"
	"simple_bank/internal/database"
	"simple_bank/internal/models"
	"simple_bank/internal/repository"
	"simple_bank/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	err = database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories and services
	db := database.GetDB()
	repo := repository.NewRepository(db)

	// Initialize services
	accountService := service.NewAccountService(repo)
	transferService := service.NewTransferService(repo, db)

	// Create HTTP server with Gin
	router := setupRouter(accountService, transferService)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.AppPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give server time to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupRouter(accountService service.AccountService, transferService service.TransferService) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"time":   time.Now().Unix(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		handler.NewAccountHandler(api, accountService, transferService)
	}

	// Debug endpoint to list all accounts
	router.GET("/debug/accounts", func(c *gin.Context) {
		accounts, err := accountService.ListAccounts(c.Request.Context(), 1, 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"accounts": accounts})
	})

	// Debug endpoint to check database
	router.GET("/debug/db", func(c *gin.Context) {
		db := database.GetDB()
		var accountCount, entryCount, transferCount int64

		db.Model(&models.Account{}).Count(&accountCount)
		db.Model(&models.Entry{}).Count(&entryCount)
		db.Model(&models.Transfer{}).Count(&transferCount)

		c.JSON(http.StatusOK, gin.H{
			"accounts":  accountCount,
			"entries":   entryCount,
			"transfers": transferCount,
		})
	})

	return router
}
