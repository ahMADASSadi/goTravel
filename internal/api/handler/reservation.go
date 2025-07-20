package handler

import (
	"encoding/json"
	"time"

	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/ahMADASSadi/goTravel/internal/models"
	response "github.com/ahMADASSadi/goTravel/internal/responses"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func CreateReservation(c *gin.Context) {
	var input models.CreateReservationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	// Decode search ID
	travelID, busID, err := services.DecodeSearchID(input.SearchID)
	if err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	// Start transaction
	tx := db.DB.Begin()
	if tx.Error != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	// Check seat availability
	var count int64
	if err := tx.Model(&models.Seat{}).
		Where("bus_id = ? AND number IN ?", busID, input.SeatNo).
		Where("available = ?", true).
		Count(&count).Error; err != nil {
		tx.Rollback()
		response.Error(c, errors.ErrServerError)
		return
	}
	if count != int64(len(input.SeatNo)) {
		tx.Rollback()
		response.Error(c, errors.ErrSeatsUnavailable)
		return
	}

	// Mark seats as unavailable
	if err := tx.Model(&models.Seat{}).
		Where("bus_id = ? AND number IN ?", busID, input.SeatNo).
		Update("available", false).Error; err != nil {
		tx.Rollback()
		response.Error(c, errors.ErrServerError)
		return
	}

	// Create reservation
	seatsJSON, _ := json.Marshal(input.SeatNo)
	reservation := models.Reservation{
		Seats:     datatypes.JSON(seatsJSON),
		TravelID:  travelID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := tx.Create(&reservation).Error; err != nil {
		tx.Rollback()
		response.Error(c, errors.ErrNotCreated)
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	response.Success(c, gin.H{
		"reservation_id": reservation.ID,
	})
}

func CancelReservation(c *gin.Context) {
	var input models.CancelReservationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	if err := db.DB.Model(&models.Reservation{}).
		Where("id = ?", input.ReservationID).
		Updates(map[string]interface{}{
			"expires_at": nil,
			"status":     models.Canceled,
		}).Error; err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}
	response.Success(c, gin.H{"message": "Reservation canceled"})
}
