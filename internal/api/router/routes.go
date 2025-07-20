package router

import (
	"github.com/ahMADASSadi/goTravel/internal/api/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(server *gin.Engine) {
	v1 := server.Group("/api/v1")
	ReservationRoutes(v1)
	OriginRoutes(v1)
	TravelRouter(v1)
}

func ReservationRoutes(group *gin.RouterGroup) {
	group.Group("/reservation").POST("/", handler.CreateReservation)
	group.Group("/reservation").POST("/cancel", handler.CancelReservation)
}

func OriginRoutes(group *gin.RouterGroup) {
	group.Group("/origins").POST("/", handler.GetOriginHandler)
	group.Group("/origins").GET("/", handler.GetOriginsHandler)
	group.Group("/destinations").POST("/", handler.GetDestinationHandler)
	group.Group("/destinations").GET("/", handler.GetDestinationHandler)
	group.Group("/terminals").POST("/", handler.GetTerminalsHandler)
}

func TravelRouter(group *gin.RouterGroup) {
	group.Group("/travels").GET("/", handler.SearchForTravel)

}
