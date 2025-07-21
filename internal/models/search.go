package models

import "time"

type TravelSearchInput struct {
	Origin        string    `form:"origin" binding:"required"`
	Destination   string    `form:"destination" binding:"required"`
	DepartureDate time.Time `form:"departure_date" time_format:"2006-01-02"`
}

type TravelSearch struct {
	WeeklyTimeSchedule
	ApproxArrivalTime time.Time `json:"approx_arrival_time" example:"2024-07-21T20:00:00Z"`
	SearchID          string    `json:"search_id" example:"123e4567-e89b-12d3-a456-426614174000"`
}
