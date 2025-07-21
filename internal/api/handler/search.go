package handler

import (
	"time"

	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/ahMADASSadi/goTravel/internal/models"
	response "github.com/ahMADASSadi/goTravel/internal/responses"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchForTravel godoc
// @Summary      Search for available travels
// @Description  Returns a list of weekly time schedules for travels that match origin, destination, and date
// @Tags         Travels
// @Accept       json
// @Produce      json
// @Param        origin          query     string  true  "Origin city code"       example: ISF
// @Param        destination     query     string  true  "Destination city code"  example: THR
// @Param        departure_date  query     string  false  "Departure date (YYYY-MM-DD)"  example: 2024-07-21
// @Success      200  {object}   map[string][]TravelSearchOutput  "List of matching travel schedules"
// @Failure      400  {object}   map[string]string  "Invalid query params"
// @Failure      500  {object}   map[string]string  "Failed to fetch data"
// @Router       /api/v1/travels [get]
func SearchForTravel(c *gin.Context) {
	var results []models.WeeklyTimeSchedule
	var input models.TravelSearchInput
	if err := c.ShouldBindQuery(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	dayOfWeek := services.SetDayOfWeek(&input.DepartureDate)

	if err := db.DB.
		Preload("Bus", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Seats", "available = ?", true) // preload seats where available = true
		}).
		Where("origin_city_code = ?", input.Origin).
		Where("destination_city_code = ?", input.Destination).
		Where("day_of_week = ?", dayOfWeek).
		Find(&results).Error; err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	var outputs []models.TravelSearch
	for _, schedule := range results {
		searchID := services.EncodeSearchID(schedule.ID, schedule.BusID)
		outputs = append(outputs, models.TravelSearch{
			WeeklyTimeSchedule: schedule,
			ApproxArrivalTime:  schedule.DepartureDate.Add(time.Duration(schedule.DepartureTime) * time.Second),
			SearchID:           searchID,
		})
	}
	response.Success(c, gin.H{"travels": outputs})
}
