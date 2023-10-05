package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/gaterun/api/admin"
)

func ServeApi() {
	r := gin.Default()

	adminApi := r.Group("/admin")
	admin.AddRoutes(adminApi)

	r.Run("127.0.0.1:8001")
}
