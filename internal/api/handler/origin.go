package handler

import (
	"fmt"
	"net/http"

	"github.com/ahMADASSadi/goTravel/internal/repository"
	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	CityCode string `json:"city_code" binding:"required"`
}

func GetOrigin(c *gin.Context) {

	origins, err := repository.GetAvailableOrigins()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"origins": origins})
}

func GetTarget(c *gin.Context) {
	var origin RequestBody

	// Parse JSON body
	if err := c.ShouldBindJSON(&origin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	originCode := origin.CityCode
	if originCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city_code is required"})
		return
	}

	fmt.Printf("originCode: %v\n", originCode) // Debugging

	targets, err := repository.GetAvailableTargets(originCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"targets": targets})
}

// GetTerminals godoc
// @Summary      Get terminals for a city
// @Description  Returns a list of terminals for the specified city code if there are any routes available
// @Tags         Terminals
// @Accept       json
// @Produce      json
// @Param        origin  body      RequestBody  true  "City Code"
// @Success      200     {object}  map[string][]repository.CityTerminal
// @Failure      400     {object}  map[string]string  "Invalid JSON body or missing city_code"
// @Failure      404     {object}  map[string]string  "No terminals found for this city"
// @Router       /api/v1/terminals/ [post]
func GetTerminals(c *gin.Context) {
	var origin RequestBody

	// Parse JSON body
	if err := c.ShouldBindJSON(&origin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	originCode := origin.CityCode
	if originCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "city_code is required"})
		return
	}

	fmt.Printf("originCode: %v\n", originCode) // Debugging

	targets, err := repository.GetCityTerminals(originCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"targets": targets})
}
