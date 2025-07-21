package router

import (
	"github.com/ahMADASSadi/goTravel/internal/api/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the API route groups and registers all handlers.
// It sets up routes for reservations, tickets, origins, and travels.
func SetupRouter(server *gin.Engine) {
	v1 := server.Group("/api/v1")
	{
		registerReservationRoutes(v1)
		registerTicketRoutes(v1)
		registerOriginRoutes(v1)
		registerTravelRoutes(v1)
	}
}

// registerReservationRoutes sets up routes for reservation operations.
// Includes endpoints for creating and canceling reservations.
func registerReservationRoutes(rg *gin.RouterGroup) {
	reservation := rg.Group("/reservations")
	{
		reservation.POST("/", handler.CreateReservationHandler)
		reservation.POST("/cancel", handler.CancelReservationHandler)
	}
}

// registerTicketRoutes sets up routes for ticket operations.
// Includes endpoints for buying, refunding, and inquiring about tickets.
func registerTicketRoutes(rg *gin.RouterGroup) {
	tickets := rg.Group("/tickets")
	{
		tickets.POST("/buy", handler.CreateTicketHandler)
		tickets.POST("/refund", handler.RefundTicketHandler)
		tickets.POST("/inquiry", handler.InquiryTicketHandler)
	}
}

// registerOriginRoutes sets up routes for origin-related operations.
// Includes endpoints for retrieving origins, destinations, and terminals.
func registerOriginRoutes(rg *gin.RouterGroup) {
	origins := rg.Group("/origins")
	{
		origins.GET("/", handler.GetOriginsHandler)
		origins.POST("/", handler.GetOriginHandler)
		origins.POST("/destinations", handler.GetDestinationHandler)
		origins.GET("/destinations", handler.GetDestinationHandler)
		origins.POST("/terminals", handler.GetTerminalsHandler)
		origins.GET("/terminals", handler.GetTerminalsHandler)
	}
}

// registerTravelRoutes sets up routes for travel search operations.
// Includes an endpoint for searching available travels.
func registerTravelRoutes(rg *gin.RouterGroup) {
	travels := rg.Group("/travels")
	{
		travels.GET("/", handler.SearchForTravel)
	}
}
