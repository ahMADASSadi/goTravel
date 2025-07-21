package repository_test

import (
	"testing"
	"time"

	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/repository"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"github.com/ahMADASSadi/goTravel/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestGetTicketByNumber(t *testing.T) {
	db := testutils.SetupTestDB(t)

	// Seed a ticket
	ticket := models.Ticket{
		TicketNo:  "ABC123",
		Passenger: models.Passenger{FirstName: "John", LastName: "Doe"},
	}
	assert.NoError(t, db.Create(&ticket).Error)

	// Act
	result, err := repository.GetTicketByNumber(db, "ABC123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ABC123", result.TicketNo)

	// Negative case: ticket does not exist
	result, err = repository.GetTicketByNumber(db, "NOT_FOUND")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTravelInfoWithSearchID(t *testing.T) {
	db := testutils.SetupTestDB(t)

	// Seed bus and schedule
	bus := models.Bus{Type: "VIP", TotalSeats: 40, RemainingSeats: 20}
	assert.NoError(t, db.Create(&bus).Error)

	schedule := models.WeeklyTimeSchedule{
		OriginCityCode:          "TEH",
		OriginTerminalCode:      "THR1",
		DestinationCityCode:     "ISF",
		DestinationTerminalCode: "ISF1",
		DayOfWeek:               2,
		DepartureTime:           900,
		DepartureDate:           time.Now(),
		ApproxDurationMins:      300,
		Price:                   500000,
		BusID:                   bus.ID,
	}
	assert.NoError(t, db.Create(&schedule).Error)

	// Generate searchID
	searchID := services.EncodeSearchID(schedule.ID, bus.ID)

	// Act
	result, err := repository.GetTravelInfoWithSearchID(db, searchID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, schedule.ID, result.ID)
	assert.Equal(t, bus.ID, result.Bus.ID)

	// Negative case: Invalid search ID
	invalidSearchID := "invalidBase64"
	result, err = repository.GetTravelInfoWithSearchID(db, invalidSearchID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetReservationInfo(t *testing.T) {
	db := testutils.SetupTestDB(t)

	// Seed a reservation
	seatsJSON := []byte(`["1A", "1B"]`) // example JSON for seats
	reservation := models.Reservation{
		Seats:     datatypes.JSON(seatsJSON),
		SearchID:  "abc123",
		Status:    0, // Default is fine, but make explicit
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	assert.NoError(t, db.Create(&reservation).Error)

	// Act: fetch by ID
	result, err := repository.GetReservationInfo(db, reservation.ID)

	// Assert: reservation found
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, reservation.ID, result.ID)
	assert.JSONEq(t, string(seatsJSON), string(result.Seats), "Seats JSON should match")
	assert.Equal(t, reservation.SearchID, result.SearchID)

	// Negative case: Reservation does not exist
	result, err = repository.GetReservationInfo(db, 999)
	assert.Error(t, err)
	assert.Nil(t, result)
}
