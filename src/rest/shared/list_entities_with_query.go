package shared

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ListWithQuery[T any](c *gin.Context, findFunc func(limit, offset int, query map[string]string) ([]T, error)) {
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

	query := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if key != "limit" && key != "offset" {
			query[key] = values[0] // first value only
		}
	}

	items, err := findFunc(limit, offset, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list entities"})
		return
	}

	c.JSON(http.StatusOK, items)
}
