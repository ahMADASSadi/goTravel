package models

import (
	"time"

	"gorm.io/datatypes"
)

type ReservationStatus uint

func (s ReservationStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Success:
		return "Succeeded"
	case Canceled:
		return "Canceled"
	default:
		return "Unknown"
	}
}

const (
	Pending ReservationStatus = iota
	Success
	Canceled
)

type CreateReservationInput struct {
	SearchID    string `json:"search_id" form:"search_id" binding:"required"`
	SeatNo      []uint `json:"seat_number" form:"seat_number" binding:"required"`
	PassengerNo uint   `json:"passenger_number" form:"passenger_number" binding:"required"`
}

type CancelReservationInput struct {
	ReservationID uint `json:"reservation_id" form:"reservation_id" binding:"required"`
}

type Reservation struct {
	ID        uint              `gorm:"primaryKey" json:"id" form:"id"`
	Seats     datatypes.JSON    `json:"seats" form:"seats"`
	SearchID  string            `json:"search_id" form:"search_id"`
	Status    ReservationStatus `gorm:"not null;default:0" json:"status" form:"status"`
	CreatedAt time.Time         `json:"created_at" form:"created_at"`
	ExpiresAt time.Time         `json:"expires_at" form:"expires_at"`
}
