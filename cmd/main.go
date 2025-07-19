package main

import (
	"fmt"
	"net/http"

	_ "github.com/ahMADASSadi/goTravel/docs"

	"github.com/ahMADASSadi/goTravel/internal/api/router"
	cfg "github.com/ahMADASSadi/goTravel/internal/config"
	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	gin.SetMode("debug")
	config := cfg.LoadConfig()
	fmt.Printf("config: %v\n", config)
	db.ConnectDatabase(config.DBDriver, config.DBSource)
	server := gin.Default()
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.SetupRouter(server)
	// @Summary      Welcome Message
	// @Description  Returns a simple welcome message for health check or API root
	// @Tags         Root
	// @Produce      json
	// @Success      200  {object}  map[string]string
	// @Router       / [get]
	server.GET("", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"welcome": "welcome"}) })

	server.Run(":8000")
}
