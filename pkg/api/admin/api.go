package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/hosseintrz/gaterun/config/models"
	cfg_persistence "github.com/hosseintrz/gaterun/config/persistence"
	"github.com/hosseintrz/gaterun/pkg/api/admin/models"
	"github.com/hosseintrz/gaterun/pkg/api/admin/persistence"
	"github.com/hosseintrz/gaterun/pkg/api/util"
)

type EndpointCallBack func(conf config.ServiceConfig) error

var endpointChangeCallbacks []EndpointCallBack

func RegisterEndpointCallback(fn EndpointCallBack) {
	endpointChangeCallbacks = append(endpointChangeCallbacks, fn)
}

func RunEndpointCallbacks(conf config.ServiceConfig) error {
	for _, fn := range endpointChangeCallbacks {
		err := fn(conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddRoutes(r *gin.RouterGroup) {
	addConsumerRoutes(r.Group("/consumers"))
	addEndpointRoutes(r.Group("/endpoints"))
	addBackendRoutes(r.Group("/backends"))
	addServiceRoutes(r.Group("/service"))
}

func addServiceRoutes(r *gin.RouterGroup) {
	r.POST("/refresh", func(c *gin.Context) {
		ctx := c.Request.Context()

		conf, err := cfg_persistence.AssembleConfig(ctx)
		if err != nil {
			util.WriteError(c, err)
			return
		}

		err = RunEndpointCallbacks(*conf)
		if err != nil {
			util.WriteError(c, err)
			return
		}

		util.WriteJSON(c, util.ResponseOk(nil), nil)
		return
	})
}

func addEndpointRoutes(r *gin.RouterGroup) {
	r.POST("/", func(c *gin.Context) {
		ctx := c.Request.Context()

		var req models.EndpointRequestDTO
		if err := c.BindJSON(&req); err != nil {
			util.WriteError(c, err)
			return
		}

		id, err := persistence.InsertEndpoint(ctx, &req)
		if err != nil {
			util.WriteError(c, err)
			return
		}

		// conf, err := assembleConfig(ctx)
		// if err != nil {
		// 	return
		// }

		// err = runEndpointCallbacks(*conf)
		// if err != nil {
		// 	util.WriteError(c, err)
		// 	return
		// }

		res := util.ResponseCreated(&gin.H{
			"id": id,
		})

		util.WriteJSON(c, res, err)
	})

	r.DELETE("/:id", func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		endpoint, err := persistence.DeleteEndpoint(ctx, id)
		if err != nil {
			util.WriteError(c, err)
			return
		}

		util.WriteJSON(c, util.NewResponse(endpoint, http.StatusNoContent), err)
	})

	r.GET("/:id", func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		endpoint, err := persistence.FetchEndpoint(ctx, id)
		util.WriteJSON(c, util.NewResponse(endpoint, http.StatusOK), err)
	})

	r.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()

		endpoints, err := persistence.FetchAllEndpoints(ctx)
		util.WriteJSON(c, util.NewResponse(endpoints, http.StatusOK), err)
	})
}

func addConsumerRoutes(r *gin.RouterGroup) {
	r.POST("/", func(c *gin.Context) {
		ctx := c.Request.Context()

		var req models.Consumer
		if err := c.BindJSON(&req); err != nil {
			util.WriteError(c, err)
			return
		}

		id, err := persistence.InsertConsumer(ctx, &req)
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

		consumer, err := persistence.GetConsumer(ctx, id)
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

		consumer, err := persistence.DeleteConsumer(ctx, id)
		util.WriteJSON(c, util.NewResponse(consumer, http.StatusNoContent), err)
	})

	r.PATCH("/:id", func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		var consumer models.Consumer
		if err := c.BindJSON(&consumer); err != nil {
			util.WriteError(c, err)
			return
		}

		res, err := persistence.UpdateConsumer(ctx, id, &consumer)
		util.WriteJSON(c, util.ResponseOk(res), err)
	})

	r.POST("/:id/key-auth", func(c *gin.Context) {
		ctx := c.Request.Context()

		consumerId, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		key, err := persistence.GenerateApiKey(ctx, consumerId)
		util.WriteJSON(c, util.ResponseCreated(key), err)
	})
}

func addBackendRoutes(r *gin.RouterGroup) {
	r.POST("/", func(c *gin.Context) {
		ctx := c.Request.Context()

		var req config.BackendConfig
		if err := c.BindJSON(&req); err != nil {
			util.WriteError(c, err)
			return
		}

		id, err := persistence.InsertBackend(ctx, &req)
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

		backend, err := persistence.FetchBackend(ctx, id)
		if err != nil {
			util.WriteError(c, err)
			return
		}

		util.WriteJSON(c, &util.Response{
			Data:       backend,
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

		backend, err := persistence.DeleteBackend(ctx, id)
		util.WriteJSON(c, util.NewResponse(backend, http.StatusNoContent), err)
	})

	r.PATCH("/:id", func(c *gin.Context) {
		ctx := c.Request.Context()

		id, err := util.Int64Param(c, "id")
		if err != nil {
			util.WriteError(c, util.NewHTTPError(http.StatusBadRequest, "invalid param id", err))
			return
		}

		var backend config.BackendConfig
		if err := c.BindJSON(&backend); err != nil {
			util.WriteError(c, err)
			return
		}

		res, err := persistence.UpdateBackend(ctx, id, &backend)
		util.WriteJSON(c, util.ResponseOk(res), err)
	})
}
