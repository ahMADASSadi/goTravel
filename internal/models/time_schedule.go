package models

import "gorm.io/gorm"

type Bus struct {
	BusCapacity uint   `json:"bus_capacity" gorm:"not null"`    // Number of seats
	BusType     string `json:"bus_type" gorm:"size:5;not null"` // Chair formation code (5 letters)
}

type WeeklyTimeSchedule struct {
	gorm.Model

	Origin         string `json:"origin" gorm:"size:5;not null"`
	OriginTerminal string `json:"origin_terminal" gorm:"size:5;not null"`

	Target         string `json:"target" gorm:"size:5;not null"`
	TargetTerminal string `json:"target_terminal" gorm:"size:5;not null"`

	TravelDay      uint8  `json:"travel_day" gorm:"not null"`
	TravelTime     uint32 `json:"travel_time" gorm:"not null"`     // Seconds since midnight (e.g., 3600 = 01:00 AM)
	TravelDuration uint16 `json:"travel_duration" gorm:"not null"` // Duration in minutes

	// Embed Bus fields directly
	Bus

	TicketPrice uint16 `json:"ticket_price" gorm:"not null"` // Price in smallest currency unit (e.g., cents)
}
