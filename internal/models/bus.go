package models

import (
	"time"

	"gorm.io/gorm"
)

type Seat struct {
	ID        uint `gorm:"primaryKey"`
	BusID     uint
	Number    int  `gorm:"not null"`
	Available bool `gorm:"default:true;not null"`
}

type Bus struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Type           string `gorm:"size:5;not null"`
	TotalSeats     int    `gorm:"default:24;not null"`
	RemainingSeats int    `gorm:"default:24;not null"`

	Seats []Seat `gorm:"foreignKey:BusID;constraint:OnDelete:CASCADE;"`
}

// Also a database trigger could be implemented to do this in the database layer rather tahn the framework
func (b *Bus) AfterCreate(tx *gorm.DB) (err error) {
	seats := make([]Seat, b.TotalSeats)

	for i := 0; i < b.TotalSeats; i++ {
		seats[i] = Seat{
			BusID:     b.ID,
			Number:    i + 1,
			Available: true,
		}
	}

	if err := tx.Create(&seats).Error; err != nil {
		return err
	}

	return nil
}
func (s *Seat) AfterSave(tx *gorm.DB) error {
	if s.BusID == 0 {
		// Load BusID from DB
		if err := tx.Model(&Seat{}).
			Where("id = ?", s.ID).
			Pluck("bus_id", &s.BusID).Error; err != nil {
			return err
		}
	}

	var bus Bus
	if err := tx.First(&bus, s.BusID).Error; err != nil {
		return err
	}
	return bus.UpdateRemainingSeats(tx)
}

func (s *Seat) AfterDelete(tx *gorm.DB) error {
	var bus Bus
	if err := tx.First(&bus, s.BusID).Error; err != nil {
		return err
	}
	return bus.UpdateRemainingSeats(tx)
}

func (b *Bus) UpdateRemainingSeats(db *gorm.DB) error {
	var count int64
	err := db.Model(&Seat{}).
		Where("bus_id = ? AND available = ?", b.ID, true).
		Count(&count).Error
	if err != nil {
		return err
	}

	return db.Model(b).Update("remaining_seats", count).Error
}
