package repository

import (
	"fmt"

	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/models"
)

func FilterAvailableOrigins(names []string) ([]string, error) {
	var filtered []string
	err := db.DB.Model(&models.WeeklyTimeSchedule{}).
		Joins("JOIN buses ON weekly_time_schedules.bus_id = buses.id").
		Where("weekly_time_schedules.origin_city_code IN ? AND buses.remaining_seats > 0", names).
		Distinct("weekly_time_schedules.origin_city_code").
		Pluck("weekly_time_schedules.origin_city_code", &filtered).Error

	return filtered, err
}

func GetAvailableOrigins() ([]string, error) {
	var originCities []string

	err := db.DB.Model(&models.WeeklyTimeSchedule{}).
		Joins("JOIN buses ON weekly_time_schedules.bus_id = buses.id").
		Where("buses.remaining_seats > 0").
		Distinct("weekly_time_schedules.origin_city_code").
		Pluck("weekly_time_schedules.origin_city_code", &originCities).Error

	return originCities, err
}

func GetAvailableDestinations(origin string) ([]string, error) {
	var targets []string
	err := db.DB.Model(&models.WeeklyTimeSchedule{}).
		Joins("JOIN buses ON weekly_time_schedules.bus_id = buses.id").
		Where("weekly_time_schedules.origin_city_code = ? AND buses.remaining_seats > 0", origin).
		Distinct("weekly_time_schedules.destination_city_code").
		Pluck("weekly_time_schedules.destination_city_code", &targets).Error

	return targets, err
}

// func GetCityTerminals(origin string) (models.CityTerminal, error) {
// 	var terminalNames []string

// 	err := db.DB.Model(&models.WeeklyTimeSchedule{}).
// 		Joins("JOIN buses ON weekly_time_schedules.bus_id = buses.id").
// 		Where("weekly_time_schedules.origin_city_code = ? AND buses.remaining_seats > 0", origin).
// 		Distinct("weekly_time_schedules.origin_terminal_code").
// 		Pluck("weekly_time_schedules.origin_terminal_code", &terminalNames).Error

// 	if err != nil {
// 		return models.CityTerminal{}, err
// 	}

// 	terminals := make([]models.Terminal, len(terminalNames))
// 	for i, name := range terminalNames {
// 		terminals[i] = models.Terminal{TerminalName: name}
// 	}

// 	return models.CityTerminal{
// 		CityName:  origin,
// 		Terminals: terminals,
// 	}, nil
// }

func GetCityTerminals(cityCode string, isOrigin bool) (models.CityTerminal, error) {
	var terminalNames []string
	var column, cityColumn string

	if isOrigin {
		column = "origin_terminal_code"
		cityColumn = "origin_city_code"
	} else {
		column = "destination_terminal_code"
		cityColumn = "destination_city_code"
	}

	err := db.DB.Model(&models.WeeklyTimeSchedule{}).
		Joins("JOIN buses ON weekly_time_schedules.bus_id = buses.id").
		Where(fmt.Sprintf("weekly_time_schedules.%s = ? AND buses.remaining_seats > 0", cityColumn), cityCode).
		Distinct(fmt.Sprintf("weekly_time_schedules.%s", column)).
		Pluck(fmt.Sprintf("weekly_time_schedules.%s", column), &terminalNames).Error

	if err != nil {
		return models.CityTerminal{}, err
	}

	terminals := make([]models.Terminal, len(terminalNames))
	for i, name := range terminalNames {
		terminals[i] = models.Terminal{TerminalName: name}
	}

	return models.CityTerminal{
		CityName:  cityCode,
		Terminals: terminals,
	}, nil
}
