package handler

import (
	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/ahMADASSadi/goTravel/internal/models"
	response "github.com/ahMADASSadi/goTravel/internal/responses"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateReservationHandler godoc
// @Summary      Create a new reservation
// @Description  Books seats for a specific travel using the provided SearchID and returns the reservation ID.
// @Tags         Reservations
// @Accept       json
// @Produce      json
// @Param        reservation  body      models.CreateReservationInput  true  "Reservation input data"
// @Success      200          {object}  map[string]uint                "Reservation successfully created"
// @Failure      400          {object}  map[string]string              "Bad request (invalid input)"
// @Failure      409          {object}  map[string]string              "Seats unavailable"
// @Failure      500          {object}  map[string]string              "Internal server error"
// @Router       /api/v1/reservations/ [post]
func CreateReservationHandler(c *gin.Context) {
	var input models.CreateReservationInput
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	if len(input.SeatNo) > int(input.PassengerNo) {
		response.Error(c, errors.ErrSeatLessThanPassenger)
		return
	}
	_, busID, err := services.DecodeSearchID(input.SearchID)
	if err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	reservation, err := services.MarkSeatsAsReserved(tx, busID, input.SearchID, input.SeatNo)
	if err != nil {
		tx.Rollback()
		response.Error(c, errors.ErrSeatsUnavailable)

		return
	}

	if err := tx.Commit().Error; err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	response.Success(c, gin.H{
		"reservation_id": reservation.ID,
	})
}

// CancelReservationHandler godoc
// @Summary      Cancel a reservation
// @Description  Cancels a reservation by ID and releases the reserved seats.
// @Tags         Reservations
// @Accept       json
// @Produce      json
// @Param        reservation  body      models.CancelReservationInput  true  "Reservation ID to cancel"
// @Success      200          {object}  map[string]string               "Reservation successfully canceled"
// @Failure      400          {object}  map[string]string               "Bad request (invalid input)"
// @Failure      404          {object}  map[string]string               "Reservation not found"
// @Failure      500          {object}  map[string]string               "Internal server error"
// @Router       /api/v1/reservations/cancel/ [post]
func CancelReservationHandler(c *gin.Context) {
	var input models.CancelReservationInput
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	if err := services.UnmarkReservedSeats(tx, input.ReservationID); err != nil {
		tx.Rollback()
		response.Error(c, errors.ErrNotFound)

		return
	}

	if err := tx.Commit().Error; err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	response.Success(c, gin.H{"message": "Reservation canceled"})
}
