package shared

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ListEntities[T any](c *gin.Context, findFunc func(limit, offset int) ([]T, error)) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	items, err := findFunc(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list entities"})
		return
	}

	c.JSON(http.StatusOK, items)
}
