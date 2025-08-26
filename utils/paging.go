package utils

import (
	"errors"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Paging(c *gin.Context) (limit int, offset int, err error) {
	limit = 100
	offset = 0
	if l, exist := c.GetQuery("limit"); exist {
		limit, err = strconv.Atoi(l)
		if err != nil {
			return 0, 0, errors.New("cannot parse limit value")
		}
	}
	if of, exist := c.GetQuery("offset"); exist {
		offset, err = strconv.Atoi(of)
		if err != nil {
			return 0, 0, errors.New("cannot parse offset value")
		}
	}
	limit = int(math.Min(float64(limit), 100))
	return limit, offset, err
}
