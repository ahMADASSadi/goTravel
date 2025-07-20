package models

import (
	"time"

	"gorm.io/gorm"
)

type WeeklyTimeSchedule struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	OriginCityCode          string `gorm:"size:5;not null" json:"origin_city_code"`
	OriginTerminalCode      string `gorm:"size:5;not null" json:"origin_city_terminal"`
	DestinationCityCode     string `gorm:"size:5;not null" json:"destination_city_code"`
	DestinationTerminalCode string `gorm:"size:5;not null" json:"destination_city_terminal"`

	DayOfWeek          uint      `gorm:"not null" json:"day_of_week"`
	DepartureTime      int64     `gorm:"not null" json:"departure_time"`
	DepartureDate      time.Time `gorm:"not null" json:"departure_date"`
	ApproxDurationMins uint      `gorm:"not null" json:"approx_duration_mins"`
	Price              uint      `gorm:"not null" json:"price"`

	BusID uint `gorm:"not null" json:"bus_id"`
	Bus   Bus  `gorm:"foreignKey:BusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (w *WeeklyTimeSchedule) BeforeSave(tx *gorm.DB) (err error) {
	if !w.DepartureDate.IsZero() {
		day := int(w.DepartureDate.Weekday())
		day = (day + 1) % 7

		w.DayOfWeek = uint(day)
	}
	return nil
}
