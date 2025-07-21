package repository

import (
	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/models"
)

func GetTicketByNumber(ticketNumber string) (*models.Ticket, error) {
	var ticket models.Ticket

	if err := db.DB.Where("ticket_no = ?", ticketNumber).First(&ticket).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

// func GetTravelInfoWithSearchID(searchID string) (*models.TravelSearch, error) {
// 	var travelSearch models.TravelSearch
// 	travelInfo, _, err := services.DecodeSearchID(searchID)
// }
