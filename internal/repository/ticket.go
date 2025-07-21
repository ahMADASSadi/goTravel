package repository

import (
	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"gorm.io/gorm"
)

func GetTicketByNumber(db *gorm.DB, ticketNumber string) (*models.Ticket, error) {
	var ticket models.Ticket

	if err := db.Where("ticket_no = ?", ticketNumber).First(&ticket).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func GetTravelInfoWithSearchID(db *gorm.DB, searchID string) (*models.WeeklyTimeSchedule, error) {
	var travel models.WeeklyTimeSchedule
	travelID, _, err := services.DecodeSearchID(searchID)
	if err != nil {
		return nil, err
	}
	if err := db.
		Where("id = ?", travelID).
		Preload("Bus").
		Preload("Bus.Seats").
		First(&travel).Error; err != nil {
		return nil, err
	}
	return &travel, nil
}

func GetReservationInfo(db *gorm.DB, reservationID uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := db.Where("id = ?", reservationID).First(&reservation).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}
