package models

import "time"

type Passenger struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	SocialCode  string `json:"social_code" binding:"required,len=10,numeric"`
}

type Ticket struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketNo  string    `gorm:"unique" json:"ticket_number"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
