package routes

import (
	"net/http"
	"simple_bank/server/internal/handler"
	"simple_bank/server/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(services *services.Services) *gin.Engine {
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
		handler.NewServicesHandler(api, services)
	}

	return router
}
