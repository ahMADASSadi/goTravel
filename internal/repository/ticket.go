package repository

import (
	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/services"
)

func GetTicketByNumber(ticketNumber string) (*models.Ticket, error) {
	var ticket models.Ticket

	if err := db.DB.Where("ticket_no = ?", ticketNumber).First(&ticket).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func GetTravelInfoWithSearchID(searchID string) (*models.WeeklyTimeSchedule, error) {
	var travel models.WeeklyTimeSchedule
	travelID, _, err := services.DecodeSearchID(searchID)
	if err != nil {
		return nil, err
	}
	if err := db.DB.
		Where("id = ?", travelID).
		Preload("Bus").
		Preload("Bus.Seats").
		First(&travel).Error; err != nil {
		return nil, err
	}
	return &travel, nil
}

func GetReservationInfo(reservationID uint) (*models.Reservation, error) {
	var reservation models.Reservation
	if err := db.DB.Where("id = ?", reservationID).First(&reservation).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}
