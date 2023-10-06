package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hosseintrz/gaterun/pkg/api/admin"
	"github.com/sirupsen/logrus"
)

func ServeApi() {
	r := gin.Default()

	adminApi := r.Group("/admin")
	admin.AddRoutes(adminApi)

	addr := "127.0.0.1:8001"
	logrus.Infof("serving admin api on %s", addr)
	r.Run(addr)
}
