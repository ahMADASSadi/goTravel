package main

import (
	"context"
	"fmt"
	"log"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go services.AutoCancelReservations(ctx)
	if err := server.Run(":8000"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
