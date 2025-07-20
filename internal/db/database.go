package db

import (
	"log"

	"github.com/ahMADASSadi/goTravel/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dbDriver, dbSource string) {
	var err error
	if dbDriver == "sqlite3" {
		DB, err = gorm.Open(sqlite.Open(dbSource), &gorm.Config{})
	}
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected")

	migrate()
}

func migrate() {
	// Auto-migrate your tables here
	err := DB.AutoMigrate(
		// &Origin{},
		// &Target{},
		// &Terminal{},
		// &Reservation{},
		// &Purchase{},
		&models.Reservation{},
		&models.Seat{},
		&models.Bus{},
		&models.WeeklyTimeSchedule{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	log.Println("Database migrated")
}
