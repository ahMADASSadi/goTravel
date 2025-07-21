package models

import (
	"time"

	"gorm.io/gorm"
)

type WeeklyTimeSchedule struct {
	ID        uint      `gorm:"primarykey" form:"id" json:"id"`
	CreatedAt time.Time `form:"created_at" json:"created_at"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at"`

	OriginCityCode          string `gorm:"size:5;not null" form:"origin_city_code" json:"origin_city_code"`
	OriginTerminalCode      string `gorm:"size:5;not null" form:"origin_city_terminal" json:"origin_city_terminal"`
	DestinationCityCode     string `gorm:"size:5;not null" form:"destination_city_code" json:"destination_city_code"`
	DestinationTerminalCode string `gorm:"size:5;not null" form:"destination_city_terminal" json:"destination_city_terminal"`

	DayOfWeek          uint      `gorm:"not null" form:"day_of_week" json:"day_of_week"`
	DepartureTime      int64     `gorm:"not null" form:"departure_time" json:"departure_time"`
	DepartureDate      time.Time `gorm:"not null" form:"departure_date" json:"departure_date"`
	ApproxDurationMins uint      `gorm:"not null" form:"approx_duration_mins" json:"approx_duration_mins"`
	Price              uint      `gorm:"not null" form:"price" json:"price"`

	BusID uint `gorm:"not null" form:"bus_id" json:"bus_id"`
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
