package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/gaterun/api/util"
)

func AddRoutes(r *gin.RouterGroup) {
	addConsumerRoutes(r.Group("/consumers"))
}

func addConsumerRoutes(r *gin.RouterGroup) {
	r.POST("/", func(c *gin.Context) {
		ctx := c.Request.Context()

		var req Consumer
		if err := c.BindJSON(&req); err != nil {
			util.WriteError(c, err)
			return
		}

		id, err := insertConsumer(ctx, &req)
		res := util.ResponseCreated(
			&gin.H{
				"id": id,
			},
		)

		util.WriteJSON(c, res, err)
	})

	r.GET("/:id", func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		consumer, err := getConsumer(ctx, id)
		if err != nil {
			util.WriteError(c, err)
			return
		}

		util.WriteJSON(c, &util.Response{
			Data:       consumer,
			StatusCode: http.StatusOK,
		}, err)
	})

	r.DELETE("/:id", func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		consumer, err := deleteConsumer(ctx, id)
		util.WriteJSON(c, util.NewResponse(consumer, http.StatusNoContent), err)
	})

	r.PATCH("/:id", func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		var consumer Consumer
		if err := c.BindJSON(&consumer); err != nil {
			util.WriteError(c, err)
			return
		}

		res, err := updateConsumer(ctx, id, &consumer)
		util.WriteJSON(c, util.ResponseOk(res), err)
	})
}
