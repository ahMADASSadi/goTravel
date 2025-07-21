package services_test

import (
	"testing"
	"time"

	"github.com/ahMADASSadi/goTravel/internal/services"
)

func TestSetDayOfWeek(t *testing.T) {
	// Test when date is nil (should return current day)
	currentDay := int(time.Now().Weekday())
	dayOfWeek := services.SetDayOfWeek(nil)
	if dayOfWeek != (currentDay+1)%7 { // Day of week should be the next day (Sunday=0, Monday=1,...)
		t.Errorf("expected %d, got %d", (currentDay+1)%7, dayOfWeek)
	}

	// Test with a specific date (e.g., a Wednesday)
	testDate := time.Date(2025, 7, 23, 0, 0, 0, 0, time.UTC) // July 23, 2025 (Wednesday)
	dayOfWeek = services.SetDayOfWeek(&testDate)
	if dayOfWeek != 4 {
		t.Errorf("expected %d, got %d", 4, dayOfWeek)
	}

	// Test with a zeroed date (should default to current time)
	zeroedDate := time.Time{} // zeroed time
	dayOfWeek = services.SetDayOfWeek(&zeroedDate)
	if dayOfWeek != (currentDay+1)%7 { // Should default to current day
		t.Errorf("expected %d, got %d", (currentDay+1)%7, dayOfWeek)
	}
}


func TestEncodeSearchID(t *testing.T) {
    scheduleID := uint(123)
    busID := uint(456)

    encodedID := services.EncodeSearchID(scheduleID, busID)
    if len(encodedID) != 20 {
        t.Errorf("expected encoded ID length of 20, got %d", len(encodedID))
    }

    // Check if the ID contains the scheduleID and busID
    decodedScheduleID, decodedBusID, err := services.DecodeSearchID(encodedID)
    if err != nil {
        t.Fatalf("failed to decode search ID: %v", err)
    }

    if decodedScheduleID != scheduleID {
        t.Errorf("expected scheduleID %d, got %d", scheduleID, decodedScheduleID)
    }

    if decodedBusID != busID {
        t.Errorf("expected busID %d, got %d", busID, decodedBusID)
    }
}


func TestDecodeSearchID(t *testing.T) {
    scheduleID := uint(123)
    busID := uint(456)

    // Encode a known scheduleID and busID
    encodedID := services.EncodeSearchID(scheduleID, busID)

    // Decode the encoded ID
    decodedScheduleID, decodedBusID, err := services.DecodeSearchID(encodedID)
    if err != nil {
        t.Fatalf("failed to decode search ID: %v", err)
    }

    if decodedScheduleID != scheduleID {
        t.Errorf("expected scheduleID %d, got %d", scheduleID, decodedScheduleID)
    }

    if decodedBusID != busID {
        t.Errorf("expected busID %d, got %d", busID, decodedBusID)
    }

    // Test invalid searchID (e.g., too short or malformed)
    invalidSearchID := "invalid_search_id"
    _, _, err = services.DecodeSearchID(invalidSearchID)
    if err == nil {
        t.Fatal("expected error for invalid search ID, but got nil")
    }
}
