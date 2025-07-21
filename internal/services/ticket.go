package services

import (
	"fmt"
	"time"

	"github.com/ahMADASSadi/goTravel/internal/models"
	"gorm.io/gorm"
)

func GenerateTicketNo(index int) string {
	return fmt.Sprintf("TKT-%d%d", index, time.Now().UnixNano())
}
func GetSearchIDFromReservation(db *gorm.DB, reservationID uint) (string, error) {
	var searchID string
	if err := db.Model(&models.Reservation{}).
		Select("search_id").
		Where("id = ?", reservationID).
		Scan(&searchID).Error; err != nil {
		return "", err
	}

	// Explicitly return an error if no searchID was found
	if searchID == "" {
		return "", fmt.Errorf("reservation not found")
	}

	return searchID, nil
}

func CreateTicket(tx *gorm.DB, ticketInput *models.TicketCreateInput) (*[]models.Ticket, error) {
	reservationID := ticketInput.ReservationID

	searchID, err := GetSearchIDFromReservation(tx, reservationID)
	if err != nil {
		return nil, err
	}

	var tickets []models.Ticket
	for i, v := range ticketInput.Passengers {
		ticket := models.Ticket{
			TicketNo:      GenerateTicketNo(i),
			ReservationID: reservationID,
			SearchID:      searchID,
			Passenger:     v,
			CreatedAt:     time.Now(),
		}
		tickets = append(tickets, ticket)
	}

	if err := tx.Create(&tickets).Error; err != nil {
		return nil, err
	}
	if err := SetReservationStatus(tx, reservationID, uint(models.Success)); err != nil {
		return nil, err
	}

	return &tickets, nil
}

func RefundTicket(tx *gorm.DB, ticketNo string) error {
	var ticket models.Ticket

	if err := tx.Where("ticket_no = ?", ticketNo).First(&ticket).Error; err != nil {
		return err
	}

	if err := tx.Model(&ticket).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	// if err := UnmarkReservedSeats(tx, ticket.ReservationID); err != nil {
	// 	return err
	// }

	return nil
}
