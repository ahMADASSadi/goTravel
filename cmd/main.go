package main

import (
	"fmt"

	_ "github.com/ahMADASSadi/goTravel/docs"

	"github.com/ahMADASSadi/goTravel/internal/api/router"
	cfg "github.com/ahMADASSadi/goTravel/internal/config"
	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/services"
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
	server.Run(":8000")
	go services.AutoCancelReservations()
}
