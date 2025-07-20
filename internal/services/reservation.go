package services

import (
	"time"

	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/models"
)

func AutoCancelReservations() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		db.DB.Model(&models.Reservation{}).
			Where("status = ? AND expires_at <= ?", models.Pending, time.Now()).
			Updates(map[string]interface{}{"status": models.Canceled})
	}
}
