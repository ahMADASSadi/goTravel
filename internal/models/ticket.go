// Package models provides the model representation for the DB schema and responses return schemas
package models

import "time"

// Passenger represents a passenger's personal information used for booking tickets.
//
// Fields:
//   - FirstName: Passenger's first name (required).
//   - LastName: Passenger's last name (required).
//   - PhoneNumber: Passenger's phone number in E.164 format (required).
//   - SocialCode: National identification code (required, exactly 10 digits).
type Passenger struct {
	FirstName   string `form:"first_name" json:"first_name" binding:"required"`
	LastName    string `form:"last_name" json:"last_name" binding:"required"`
	PhoneNumber string `form:"phone_number" json:"phone_number" binding:"required,e164"`
	SocialCode  string `form:"social_code" json:"social_code" binding:"required,len=10,numeric"`
}

// Ticket represents a bus ticket issued to a passenger.
//
// Fields:
//   - ID: Auto-incremented primary key.
//   - TicketNo: Unique ticket number.
//   - CreatedAt: Timestamp when the ticket was created.
//   - DeletedAt: Timestamp when the ticket was soft-deleted.
type Ticket struct {
	ID            uint   `gorm:"primaryKey" form:"id" json:"id"`
	TicketNo      string `gorm:"unique" form:"ticket_number" json:"ticket_number"`
	SearchID      string `form:"search_id" json:"search_id"`
	ReservationID uint   `form:"reservation_id" json:"reservation_id"`
	Passenger
	CreatedAt time.Time `form:"created_at" json:"created_at"`
	DeletedAt time.Time `gorm:"default:null" form:"deleted_at" json:"deleted_at"`
}

type TicketCreateInput struct {
	ReservationID uint        `form:"reservation_id" json:"reservation_id"`
	Passengers    []Passenger `form:"passengers" json:"passengers"`
}

type TicketRefundInput struct {
	TicketNo string `form:"ticket_number" json:"ticket_number"`
}

type TicketInquiry struct {
	SearchID      string `form:"search_id" json:"search_id"`
	TicketNo      string `form:"ticket_number" json:"ticket_number"`
	ReservationID uint   `form:"reservation_id" json:"reservation_id"`
}
