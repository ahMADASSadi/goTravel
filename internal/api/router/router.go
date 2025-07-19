package router

import (
	"github.com/ahMADASSadi/goTravel/internal/api/handler"
	"github.com/gin-gonic/gin"
)
func SetupRouter(server *gin.Engine) {
	v1 := server.Group("/api/v1")
	// Pass the grouped v1 router to ReservationRoutes
	ReservationRoutes(v1)
	// Add other routes if needed
	OriginRoutes(v1)
}

func ReservationRoutes(group *gin.RouterGroup) {
	// Using group to define the reservations route
	// group.Group("/reservations").GET("", listReservationsHandler) // GET /api/v1/reservations
	// group.Group("/reservations").POST("", createReservationHandler) // POST /api/v1/reservations
	// group.Group("/reservations/:id").GET("", getReservationHandler) // GET /api/v1/reservations/:id
	// group.Group("/reservations/:id").PATCH("", cancelReservationHandler) // PATCH /api/v1/reservations/:id
}

func OriginRoutes(group *gin.RouterGroup) {
	// Define origins routes under /origins path
	group.Group("/origins").POST("/", handler.GetOrigin) // POST /api/v1/origins
	group.Group("/origins").GET("/", handler.GetOrigin) // POST /api/v1/origins
	group.Group("/targets").POST("/",handler.GetTarget)
	group.Group("/terminals").POST("/",handler.GetTerminals)
	// group.Group("/origins/:code").GET("", handler.GetOriginByCode) // GET /api/v1/origins/:code
}
