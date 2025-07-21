package handler

import (
	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/repository"
	response "github.com/ahMADASSadi/goTravel/internal/responses"
	"github.com/gin-gonic/gin"
)

// GetOriginsHandler godoc
// @Summary      Get all the available origins
// @Description  Returns a list of city codes if there are any routes available
// @Tags         Cities
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string][]string  "List of available origin city codes"
// @Failure      404  {object}  map[string]string    "No available origins found"
// @Router       /api/v1/origins [get]
func GetOriginsHandler(c *gin.Context) {
	origins, err := repository.GetAvailableOrigins(db.DB)
	if err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}
	response.Success(c, gin.H{"origins": origins})
}

// GetOriginHandler godoc
// @Summary      Filter available origins
// @Description  Receives a list of city codes and returns the ones that have routes available
// @Tags         Cities
// @Accept       json
// @Produce      json
// @Param        city_codes  body      []models.RequestBody  true  "List of city codes to filter"
// @Success      200         {object}  map[string][]string  "Filtered list of available origins"
// @Failure      400         {object}  map[string]string    "Invalid JSON body or missing city_code"
// @Failure      500         {object}  map[string]string    "Internal server error"
// @Router       /api/v1/origins/ [post]
func GetOriginHandler(c *gin.Context) {
	var requests []models.RequestBody
	if err := c.ShouldBind(&requests); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}
	cityCodes := make([]string, 0, len(requests))
	for _, req := range requests {
		if req.CityCode == "" {
			response.Error(c, errors.ErrBadRequest)
			return
		}
		cityCodes = append(cityCodes, req.CityCode)
	}
	filtered, err := repository.FilterAvailableOrigins(db.DB, cityCodes)
	if err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	response.Success(c, gin.H{"available_origins": filtered})
}

// GetDestinationHandler godoc
// @Summary      Get available destinations for an origin
// @Description  Receives an origin city code and returns possible destination city codes
// @Tags         Cities
// @Accept       json
// @Produce      json
// @Param        city_code  body      models.RequestBody   true  "Origin city code"
// @Success      200        {object}  map[string][]string  "List of destination city codes"
// @Failure      400        {object}  map[string]string    "Invalid JSON body or missing city_code"
// @Failure      404        {object}  map[string]string    "No destinations found for the origin"
// @Router       /api/v1/origins/destinations/ [get]
// @Router       /api/v1/origins/destinations/ [post]
func GetDestinationHandler(c *gin.Context) {
	var origin models.RequestBody

	if err := c.ShouldBind(&origin); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	originCode := origin.CityCode
	if originCode == "" {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	targets, err := repository.GetAvailableDestinations(db.DB, originCode)
	if err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}
	response.Success(c, gin.H{"targets": targets})
}

// GetTerminalsHandler godoc
// @Summary      Get terminals for a city
// @Description  Returns a list of terminals for the specified city code. Set `is_origin` to true to get origin terminals, or false for destination terminals.
// @Tags         Cities
// @Accept       json
// @Produce      json
// @Param        terminal_request  body      models.TerminalRequestBody  true  "City code and terminal type"
// @Success      200               {object}  map[string][]models.CityTerminal  "List of terminals"
// @Failure      400               {object}  map[string]string  "Invalid JSON body or missing city_code"
// @Failure      404               {object}  map[string]string  "No terminals found for this city"  "No terminals found for this city"
// @Router       /api/v1/origins/terminals/ [get]
// @Router       /api/v1/origins/terminals/ [post]
func GetTerminalsHandler(c *gin.Context) {
	var origin models.TerminalRequestBody
	if err := c.ShouldBind(&origin); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	originCode := origin.CityCode
	if originCode == "" {
		response.Error(c, errors.ErrBadRequest)
		return
	}

	terminals, err := repository.GetCityTerminals(db.DB, originCode, origin.IsOrigin)
	if err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	response.Success(c, gin.H{"terminals": terminals})
}
