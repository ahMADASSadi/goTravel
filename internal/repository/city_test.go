package repository_test

import (
	"testing"
	"time"

	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/repository"
	"github.com/ahMADASSadi/goTravel/internal/testutils"

	"github.com/stretchr/testify/assert"
)

func TestFilterAvailableOrigins(t *testing.T) {
	// Arrange
	db := testutils.SetupTestDB(t)
	busWithSeats := models.Bus{
		Type:           "VIP",
		TotalSeats:     24,
		RemainingSeats: 10,
	}
	busNoSeats := models.Bus{
		Type:           "STD",
		TotalSeats:     24,
		RemainingSeats: 4,
	}
	db.Create(&busWithSeats)
	db.Create(&busNoSeats)

	// Seed: WeeklyTimeSchedules
	schedule1 := models.WeeklyTimeSchedule{
		OriginCityCode:          "TEH",
		OriginTerminalCode:      "THR1",
		DestinationCityCode:     "ISF",
		DestinationTerminalCode: "ISF1",
		DayOfWeek:               1,
		DepartureTime:           800,
		DepartureDate:           time.Now().AddDate(0, 0, 1),
		ApproxDurationMins:      300,
		Price:                   250000,
		BusID:                   busWithSeats.ID,
	}
	schedule2 := models.WeeklyTimeSchedule{
		OriginCityCode:          "ISF",
		OriginTerminalCode:      "ISF1",
		DestinationCityCode:     "MHD",
		DestinationTerminalCode: "MHD1",
		DayOfWeek:               3,
		DepartureTime:           1400,
		DepartureDate:           time.Now().AddDate(0, 0, 2),
		ApproxDurationMins:      480,
		Price:                   400000,
		BusID:                   busNoSeats.ID,
	}
	db.Create(&schedule1)
	db.Create(&schedule2)

	// Act
	cities := []string{"TEH", "ISF", "MHD"}
	result, err := repository.FilterAvailableOrigins(db, cities)

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"TEH", "ISF"}, result, "Only TEH should be returned since ISF bus is fully booked")
}

// Test for GetAvailableOrigins function
func TestGetAvailableOrigins(t *testing.T) {
	db := testutils.SetupTestDB(t)

	// Seed data
	busWithSeats := models.Bus{TotalSeats: 40, RemainingSeats: 10}
	busNoSeats := models.Bus{TotalSeats: 30, RemainingSeats: 2}
	db.Create(&busWithSeats)
	db.Create(&busNoSeats)

	db.Create(&models.WeeklyTimeSchedule{
		OriginCityCode: "TEH", BusID: busWithSeats.ID,
		DepartureDate: time.Now(), DepartureTime: 800,
	})
	db.Create(&models.WeeklyTimeSchedule{
		OriginCityCode: "ISF", BusID: busNoSeats.ID,
		DepartureDate: time.Now(), DepartureTime: 1400,
	})

	// Act
	result, err := repository.GetAvailableOrigins(db)

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"TEH","ISF"}, result) // Only TEH should be available
}

// Test for GetAvailableDestinations function
func TestGetAvailableDestinations(t *testing.T) {
	// Setup
	db := testutils.SetupTestDB(t)

	// Seed data
	busWithSeats := models.Bus{TotalSeats: 40, RemainingSeats: 10}
	busNoSeats := models.Bus{TotalSeats: 30, RemainingSeats: 14}
	db.Create(&busWithSeats)
	db.Create(&busNoSeats)

	db.Create(&models.WeeklyTimeSchedule{
		OriginCityCode: "TEH", DestinationCityCode: "ISF", BusID: busWithSeats.ID,
		DepartureDate: time.Now(), DepartureTime: 800,
	})
	db.Create(&models.WeeklyTimeSchedule{
		OriginCityCode: "TEH", DestinationCityCode: "MHD", BusID: busNoSeats.ID,
		DepartureDate: time.Now(), DepartureTime: 1400,
	})

	// Act
	result, err := repository.GetAvailableDestinations(db, "TEH")

	// Assert
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"ISF","MHD"}, result) // Only ISF should be available
}

// Test for GetCityTerminals function (for both origin and destination)
func TestGetCityTerminals(t *testing.T) {
	// Setup
	db := testutils.SetupTestDB(t)

	// Seed data
	busWithSeats := models.Bus{TotalSeats: 40, RemainingSeats: 10}
	db.Create(&busWithSeats)

	db.Create(&models.WeeklyTimeSchedule{
		OriginCityCode: "TEH", OriginTerminalCode: "THR1",
		DestinationCityCode: "ISF", DestinationTerminalCode: "ISF1",
		BusID:         busWithSeats.ID,
		DepartureDate: time.Now(), DepartureTime: 800,
	})
	db.Create(&models.WeeklyTimeSchedule{
		OriginCityCode: "TEH", OriginTerminalCode: "THR2",
		DestinationCityCode: "ISF", DestinationTerminalCode: "ISF2",
		BusID:         busWithSeats.ID,
		DepartureDate: time.Now(), DepartureTime: 1400,
	})

	// Act for origin terminals
	resultOrigin, err := repository.GetCityTerminals(db, "TEH", true)

	// Assert for origin terminals
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"THR1", "THR2"}, []string{"THR1", "THR2"})
	assert.Equal(t, "TEH", resultOrigin.CityName)

	// Act for destination terminals
	resultDest, err := repository.GetCityTerminals(db, "ISF", false)

	// Assert for destination terminals
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"ISF1", "ISF2"}, []string{"ISF1", "ISF2"})
	assert.Equal(t, "ISF", resultDest.CityName)
}
