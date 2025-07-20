package models

import (
	"time"

	"gorm.io/datatypes"
)

const (
	Pending uint = iota
	Succeed
	Canceled
)

type CreateReservationInput struct {
	SearchID     string `json:"search_id" binding:"required"`
	SeatNo       []uint `json:"seat_number" binding:"required"`
	PasserngerNo uint   `json:"passenger_number" binding:"required"`
}

type CancelReservationInput struct {
	ReservationID uint `json:"reservation_id" binding:"required"`
}

type Reservation struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Seats     datatypes.JSON `json:"seats"`
	TravelID  uint           `json:"travel_id"`
	Status    uint           `gorm:"not null;default:0" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `json:"expires_at"`
}
