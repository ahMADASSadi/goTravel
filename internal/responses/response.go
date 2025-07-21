package response

import (
	"net/http"

	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.ApiError:
		c.JSON(e.Code, gin.H{
			"success": false,
			"error":   e.Message,
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "unexpected server error",
		})
	}
}
