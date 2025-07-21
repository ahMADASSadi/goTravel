package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/ahMADASSadi/goTravel/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func AutoCancelReservations(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping AutoCancelReservations")
			return
		case <-ticker.C:
			db.DB.Model(&models.Reservation{}).
				Where("status = ? AND expires_at <= ?", models.Pending, time.Now()).
				Updates(map[string]interface{}{"status": models.Canceled})
		}
	}
}
func createReservation(tx *gorm.DB, seatNumbers []uint, searchID string) (*models.Reservation, error) {
	seatsJSON, _ := json.Marshal(seatNumbers)
	reservation := models.Reservation{
		Seats:     datatypes.JSON(seatsJSON),
		SearchID:  searchID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := tx.Create(&reservation).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}
func MarkSeatsAsReserved(tx *gorm.DB, busID uint, searchID string, seatNumbers []uint) (*models.Reservation, error) {
	// Check seat availability
	var count int64
	if err := tx.Model(&models.Seat{}).
		Where("bus_id = ? AND number IN ?", busID, seatNumbers).
		Where("available = ?", true).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count != int64(len(seatNumbers)) {
		return nil, errors.ErrSeatsUnavailable
	}

	// Mark seats as unavailable
	if err := tx.Model(&models.Seat{}).
		Where("bus_id = ? AND number IN ?", busID, seatNumbers).
		Update("available", false).Error; err != nil {
		return nil, err
	}

	reservation, err := createReservation(tx, seatNumbers, searchID)
	if err != nil {
		return nil, err
	}
	return reservation, nil
}
func cancelReservation(tx *gorm.DB, reservationID uint) error {
	if err := tx.Model(&models.Reservation{}).
		Where("id = ?", reservationID).
		Updates(map[string]interface{}{
			"expires_at": nil,
			"status":     models.Canceled,
		}).Error; err != nil {
		return err
	}
	return nil
}

func UnmarkReservedSeats(tx *gorm.DB, reservationID uint) error {
	var reservation models.Reservation
	if err := tx.First(&reservation, reservationID).Error; err != nil {
		return err
	}
	if reservation.Status == models.Canceled {
		return fmt.Errorf("reservation already canceled")
	}

	var seatNumbers []uint
	if err := json.Unmarshal(reservation.Seats, &seatNumbers); err != nil {
		return err
	}

	_, busID, err := DecodeSearchID(reservation.SearchID)
	if err != nil {
		return err
	}

	if err := cancelReservation(tx, reservationID); err != nil {
		return err
	}

	if err := tx.Model(&models.Seat{}).
		Where("bus_id = ? AND number IN ?", busID, seatNumbers).
		Update("available", true).Error; err != nil {
		return err
	}

	return nil
}

func SetReservationStatus(tx *gorm.DB, reservationID, status uint) error {
	if err := tx.Model(&models.Reservation{}).
		Where("id = ?", reservationID).
		Updates(map[string]interface{}{
			"expires_at": nil,
			"status":     models.Success,
		}).Error; err != nil {
		return err
	}
	return nil
}
