package services_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"github.com/ahMADASSadi/goTravel/internal/testutils"
	"gorm.io/gorm"
)

// func setupTestDB(t *testing.T) *gorm.DB {
//     db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
//     if err != nil {
//         t.Fatalf("failed to connect test db: %v", err)
//     }
//     if err := db.AutoMigrate(&models.Seat{}); err != nil {
//         t.Fatalf("failed to migrate: %v", err)
//     }
// 	if err := db.AutoMigrate(&models.Reservation{}); err != nil {
//         t.Fatalf("failed to migrate: %v", err)
//     }
//     return db
// }

func TestCreateReservation_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)

	seatNumbers := []uint{1, 2, 3}
	searchID := "test-search-id"

	tx := db.Begin()
	defer tx.Rollback()

	reservation, err := services.CreateReservation(tx, seatNumbers, searchID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify reservation data
	if reservation.SearchID != searchID {
		t.Errorf("expected SearchID %s, got %s", searchID, reservation.SearchID)
	}

	var decodedSeats []uint
	if err := json.Unmarshal(reservation.Seats, &decodedSeats); err != nil {
		t.Fatalf("failed to unmarshal seats JSON: %v", err)
	}

	if len(decodedSeats) != len(seatNumbers) {
		t.Errorf("expected %d seats, got %d", len(seatNumbers), len(decodedSeats))
	}

	// Check expiry time
	if reservation.ExpiresAt.Sub(reservation.CreatedAt) != 15*time.Minute {
		t.Errorf("expected expiry in 15 minutes, got %v", reservation.ExpiresAt.Sub(reservation.CreatedAt))
	}
}

func TestCreateReservation_Failure(t *testing.T) {
	db := testutils.SetupTestDB(t)

	seatNumbers := []uint{1, 2, 3}
	searchID := "test-search-id"

	// Use a closed transaction to force an error
	tx := db.Begin()
	tx.Rollback() // rollback immediately to make tx invalid

	_, err := services.CreateReservation(tx, seatNumbers, searchID)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestMarkSeatsAsReserved_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)

	busID := uint(1)
	seatNumbers := []uint{1, 2, 3}
	searchID := "search-123"

	// Seed seats with available seats
	seedSeats(db, busID, 10, []uint{}) // all seats available

	// Start a transaction
	tx := db.Begin()
	defer tx.Rollback() // ensures rollback after the test

	// Attempt to mark seats as reserved
	reservation, err := services.MarkSeatsAsReserved(tx, busID, searchID, seatNumbers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Ensure reservation is created correctly
	if reservation.SearchID != searchID {
		t.Errorf("expected SearchID %s, got %s", searchID, reservation.SearchID)
	}

	// Check that the seats are now marked as unavailable
	var unavailableCount int64
	if err := tx.Model(&models.Seat{}).
		Where("bus_id = ? AND number IN ? AND available = ?", busID, seatNumbers, false).
		Count(&unavailableCount).Error; err != nil {
		t.Fatalf("query error: %v", err)
	}
	if unavailableCount != int64(len(seatNumbers)) {
		t.Errorf("expected %d seats unavailable, got %d", len(seatNumbers), unavailableCount)
	}
}

// func seedSeats(db *gorm.DB, busID uint, totalSeats uint, reservedSeats []uint) {
// 	for i := uint(1); i <= totalSeats; i++ {
// 		available := true
// 		for _, r := range reservedSeats {
// 			if i == r {
// 				available = false
// 				break
// 			}
// 		}
// 		db.Create(&models.Seat{
// 			BusID:     busID,
// 			Number:    int(i), // Casting uint to int here
// 			Available: available,
// 		})
// 	}
// }

func seedSeats(db *gorm.DB, busID uint, totalSeats uint, reservedSeats []uint) {
    for i := uint(1); i <= totalSeats; i++ {
        available := true
        for _, r := range reservedSeats {
            if i == r {
                available = false
                break
            }
        }
        db.Create(&models.Seat{
            BusID:     busID,
            Number:    int(i),
            Available: available, // Make sure seat 2 is unavailable
        })
    }
}


func TestMarkSeatsAsReserved_SeatsUnavailable(t *testing.T) {
    db := testutils.SetupTestDB(t)

    busID := uint(1)
    seatNumbers := []uint{1, 2, 3}
    searchID := "search-123"

    // Pre-reserve seat 2
    seedSeats(db, busID, 10, []uint{2})

    tx := db.Begin()
    defer tx.Rollback()

    // Try to mark seats as reserved, should fail because seat 2 is unavailable
    _, err := services.MarkSeatsAsReserved(tx, busID, searchID, seatNumbers)

    // Check that an error is returned
    if err == nil {
        t.Fatal("expected error but got nil")
    }

    // Print the actual error for debugging
    fmt.Printf("Error returned: %v\n", err)

    // Check if the error is the expected "seats unavailable" error
    if err != errors.ErrSeatsUnavailable {
        t.Errorf("expected ErrSeatsUnavailable but got %v", err)
    }
}




func TestMarkSeatsAsReserved_DBError(t *testing.T) {
	db := testutils.SetupTestDB(t)

	busID := uint(1)
	seatNumbers := []uint{1, 2, 3}
	searchID := "search-123"

	seedSeats(db, busID, 10, []uint{})

	// Start and rollback transaction to simulate DB error
	tx := db.Begin()
	tx.Rollback()

	_, err := services.MarkSeatsAsReserved(tx, busID, searchID, seatNumbers)
	if err == nil {
		t.Fatal("expected DB error but got nil")
	}
}

func createTestReservation(db *gorm.DB, seatNumbers []uint, searchID string, status uint) *models.Reservation {
	seatsJSON, _ := json.Marshal(seatNumbers)
	reservation := &models.Reservation{
		Seats:    seatsJSON,
		SearchID: searchID,
		Status:   models.ReservationStatus(status),
	}
	db.Create(reservation)
	return reservation
}

func TestCancelReservation_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)

	res := createTestReservation(db, []uint{1, 2}, "bus-123", uint(models.Success))

	tx := db.Begin()
	defer tx.Rollback()

	err := services.CancelReservation(tx, res.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updated models.Reservation
	db.First(&updated, res.ID)
	if updated.Status != models.Canceled {
		t.Errorf("expected status Canceled, got %d", updated.Status)
	}
}

func TestCancelReservation_DBError(t *testing.T) {
	db := testutils.SetupTestDB(t)

	res := createTestReservation(db, []uint{1, 2}, "bus-123", uint(models.Success))

	tx := db.Begin()
	tx.Rollback() // force error

	err := services.CancelReservation(tx, res.ID)
	if err == nil {
		t.Fatal("expected DB error, got nil")
	}
}

func TestUnmarkReservedSeats_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)

	busID := uint(1)
	seatNumbers := []uint{1, 2, 3}
	reservedSeat := []uint{1, 2, 3}
	searchID := fmt.Sprintf("encoded|%d|test", busID) // Use your real EncodeSearchID
	res := createTestReservation(db, seatNumbers, searchID, uint(models.Success))
	seedSeats(db, busID, uint(len(seatNumbers)), reservedSeat)

	tx := db.Begin()
	defer tx.Rollback()

	err := services.UnmarkReservedSeats(tx, res.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check seats marked as available
	var availableCount int64
	db.Model(&models.Seat{}).
		Where("bus_id = ? AND number IN ? AND available = ?", busID, seatNumbers, true).
		Count(&availableCount)
	if availableCount != int64(len(seatNumbers)) {
		t.Errorf("expected all seats available, got %d/%d", availableCount, len(seatNumbers))
	}

	// Check reservation status
	var updated models.Reservation
	db.First(&updated, res.ID)
	if updated.Status != models.Canceled {
		t.Errorf("expected status Canceled, got %d", updated.Status)
	}
}

func TestUnmarkReservedSeats_AlreadyCanceled(t *testing.T) {
	db := testutils.SetupTestDB(t)

	busID := uint(1)
	seatNumbers := []uint{1, 2}
	searchID := fmt.Sprintf("encoded|%d|test", busID)
	res := createTestReservation(db, seatNumbers, searchID, uint(models.Canceled))

	tx := db.Begin()
	defer tx.Rollback()

	err := services.UnmarkReservedSeats(tx, res.ID)
	if err == nil {
		t.Fatal("expected error for already canceled reservation, got nil")
	}
}

func TestSetReservationStatus_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)

	res := createTestReservation(db, []uint{4, 5}, "bus-321", uint(models.Pending))

	tx := db.Begin()
	defer tx.Rollback()

	err := services.SetReservationStatus(tx, res.ID, uint(models.Success))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var updated models.Reservation
	db.First(&updated, res.ID)
	if updated.Status != models.Success {
		t.Errorf("expected status Success, got %d", updated.Status)
	}
}
