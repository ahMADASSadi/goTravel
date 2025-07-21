package router

import (
	"github.com/ahMADASSadi/goTravel/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(server *gin.Engine) {
	v1 := server.Group("/api/v1")
	{
		registerReservationRoutes(v1)
		registerTicketRoutes(v1)
		registerOriginRoutes(v1)
		registerTravelRoutes(v1)
	}
}

func registerReservationRoutes(rg *gin.RouterGroup) {
	reservation := rg.Group("/reservations")
	{
		reservation.POST("/", handler.CreateReservationHandler)
		reservation.POST("/cancel", handler.CancelReservationHandler)
	}
}

func registerTicketRoutes(rg *gin.RouterGroup) {
	tickets := rg.Group("/tickets")
	{
		tickets.POST("/buy", handler.CreateTicketHandler)
		tickets.POST("/refund", handler.RefundTicketHandler)
	}
}

func registerOriginRoutes(rg *gin.RouterGroup) {
	origins := rg.Group("/origins")
	{
		origins.GET("/", handler.GetOriginsHandler)            // list all origins
		origins.POST("/", handler.GetOriginHandler)            // filter origins
		origins.POST("/destinations", handler.GetDestinationHandler) // get destinations for origin
		origins.POST("/terminals", handler.GetTerminalsHandler)     // get terminals for city
	}
}

func registerTravelRoutes(rg *gin.RouterGroup) {
	travels := rg.Group("/travels")
	{
		travels.GET("/", handler.SearchForTravel)
	}
}