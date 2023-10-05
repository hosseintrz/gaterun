package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/gaterun/api/admin"
)

func ServeApi() {
	r := gin.Default()

	adminApi := r.Group("/admin")
	admin.AddRoutes(adminApi)

	r.Run("http://localhost:8001")
}
