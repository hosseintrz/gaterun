package util

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidParam = errors.New("param is invalid")
)

func Int64Param(c *gin.Context, name string) (int64, error) {
	str := c.Param(name)
	if len(str) == 0 {
		return 0, ErrInvalidParam
	}

	p, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return int64(p), nil
}
