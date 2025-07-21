package testutils

import (
	"testing"

	"github.com/ahMADASSadi/goTravel/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.Bus{},
		&models.Seat{},
		&models.Ticket{},
		&models.Reservation{},
		&models.WeeklyTimeSchedule{},
		// Add other models here if needed
	); err != nil {
		t.Fatalf("failed to auto-migrate models: %v", err)
	}

	return db
}
