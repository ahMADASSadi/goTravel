package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"github.com/ahMADASSadi/goTravel/internal/testutils"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db := testutils.SetupTestDB(t)
	return db
}

// Test for generateTicketNo
func TestGenerateTicketNo(t *testing.T) {
	ticketNo := services.GenerateTicketNo(1)
	assert.Contains(t, ticketNo, "TKT-1")
	assert.Len(t, ticketNo, 24) // because UnixNano is a long number
}

// Test for GetSearchIDFromReservation
func TestGetSearchIDFromReservation_Success(t *testing.T) {
	db := setupTestDB(t)

	reservation := &models.Reservation{
		SearchID: "search-123",
	}
	db.Create(reservation)

	searchID, err := services.GetSearchIDFromReservation(db, reservation.ID)
	assert.NoError(t, err)
	assert.Equal(t, "search-123", searchID)
}

func TestGetSearchIDFromReservation_Error(t *testing.T) {
	db := setupTestDB(t)

	// This reservation doesn't exist
	_, err := services.GetSearchIDFromReservation(db, 999)
	assert.Error(t, err)
}

func TestCreateTicket_Success(t *testing.T) {
	db := setupTestDB(t)

	// Generate a valid SearchID using EncodeSearchID
	scheduleID := uint(1)
	busID := uint(2)
	searchID := services.EncodeSearchID(scheduleID, busID)

	reservation := &models.Reservation{
		SearchID: searchID, // Use the generated SearchID
	}
	db.Create(reservation)

	ticketInput := &models.TicketCreateInput{
		ReservationID: reservation.ID,
		Passengers: []models.Passenger{
			{FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890", SocialCode: "1111111111"},
		},
	}

	tx := db.Begin()
	defer tx.Rollback()

	tickets, err := services.CreateTicket(tx, ticketInput)
	assert.NoError(t, err)
	assert.Len(t, *tickets, 1)
	assert.Equal(t, "TKT-0", (*tickets)[0].TicketNo[:5]) // Check ticket format
	assert.Equal(t, reservation.ID, (*tickets)[0].ReservationID)
	assert.Equal(t, searchID, (*tickets)[0].SearchID) // Match the generated SearchID
}

func TestCreateTicket_Failure_NoReservation(t *testing.T) {
	db := setupTestDB(t)

	ticketInput := &models.TicketCreateInput{
		ReservationID: 999,
		Passengers: []models.Passenger{
			{FirstName: "John", LastName: "Doe", PhoneNumber: "1234567890", SocialCode: "1111111111"},
		},
	}

	tx := db.Begin()
	defer tx.Rollback()

	tickets, err := services.CreateTicket(tx, ticketInput)
	assert.Error(t, err)
	assert.Nil(t, tickets)
}

// Test for RefundTicket
func TestRefundTicket_Success(t *testing.T) {
	db := setupTestDB(t)

	reservation := &models.Reservation{
		SearchID: "search-123",
	}
	db.Create(reservation)

	ticket := &models.Ticket{
		TicketNo:      "TKT-1",
		ReservationID: reservation.ID,
		SearchID:      reservation.SearchID,
	}
	db.Create(ticket)

	tx := db.Begin()
	defer tx.Rollback()

	err := services.RefundTicket(tx, "TKT-1")
	assert.NoError(t, err)

	var refundedTicket models.Ticket
	db.First(&refundedTicket, "ticket_no = ?", "TKT-1")
	assert.NotNil(t, refundedTicket.DeletedAt)
}

func TestRefundTicket_Failure_TicketNotFound(t *testing.T) {
	db := setupTestDB(t)

	tx := db.Begin()
	defer tx.Rollback()

	err := services.RefundTicket(tx, "TKT-NOT-EXIST")
	assert.Error(t, err)
}
