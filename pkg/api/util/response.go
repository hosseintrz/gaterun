package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data       any
	StatusCode int
}

func NewResponse(data any, status int) *Response {
	return &Response{
		Data:       data,
		StatusCode: status,
	}
}

func ResponseOk(data any) *Response {
	return NewResponse(data, http.StatusOK)
}

func ResponseCreated(data any) *Response {
	return NewResponse(data, http.StatusCreated)
}

func WriteJSON(c *gin.Context, res *Response, err error) {
	if err != nil {
		WriteError(c, err)
		return
	}

	c.JSON(res.StatusCode, res.Data)
}

func WriteError(c *gin.Context, err error) {
	if httpErr, ok := err.(*HTTPError); ok {
		c.JSON(httpErr.Status, &gin.H{
			"error": httpErr.Message,
		})
	} else {
		c.JSON(http.StatusInternalServerError, &gin.H{
			"error": err.Error(),
		})
	}
}
