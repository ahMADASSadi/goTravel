package repository

import (
	"fmt"

	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/models"
)

func GetAvailableOrigins() ([]string, error) {
	var originCities []string

	db := db.DB
	err := db.Model(&models.WeeklyTimeSchedule{}).
		Where("bus_capacity <> 0").          // Filter by origin city code
		Distinct("origin").                  // Ensure we get unique target cities
		Pluck("origin", &originCities).Error // Populate targetCities slice

	if err != nil {
		return nil, err
	}

	return originCities, nil
}

type TargetCity struct {
	Target string `json:"target"`
}

type Terminal struct {
	TerminalName string `json:"terminal_name"`
}
type CityTerminal struct {
	CityName  string     `json:"city_name"`
	Terminals []Terminal `json:"terminals"`
}

func GetAvailableTargets(origin string) ([]TargetCity, error) {
	var targets []TargetCity
	fmt.Printf("origin: %v\n", origin)

	// Query distinct targets and map to []TargetCity
	if err := db.DB.Model(&models.WeeklyTimeSchedule{}).
		Where("origin = ? AND bus_capacity <> 0", origin).
		Distinct("target").
		Pluck("target", &targets).Error; err != nil {
		return nil, err
	}

	return targets, nil
}

func GetCityTerminals(origin string) (CityTerminal, error) {
	var terminalNames []string

	// Query distinct origin terminals for this city
	if err := db.DB.Model(&models.WeeklyTimeSchedule{}).
		Where("origin = ?", origin).
		Distinct("origin_terminal").
		Pluck("origin_terminal", &terminalNames).Error; err != nil {
		return CityTerminal{}, err
	}

	// Map terminal names into []Terminal
	terminals := make([]Terminal, len(terminalNames))
	for i, name := range terminalNames {
		terminals[i] = Terminal{TerminalName: name}
	}

	// Build the CityTerminal object
	cityTerminal := CityTerminal{
		CityName:  origin,
		Terminals: terminals,
	}

	return cityTerminal, nil
}
