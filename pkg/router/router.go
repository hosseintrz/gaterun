package router

import (
	"context"
	"net/http"

	"github.com/hosseintrz/gaterun/config/models"
)

type Router interface {
	Run(models.ServiceConfig)
	Shutdown()
}

type RouterFunc func(models.ServiceConfig)

func (f RouterFunc) Run(cfg models.ServiceConfig) { f(cfg) }

type Factory interface {
	New() Router
	NewWithContext(context.Context) Router
}

type ServerRunnerFunc func(context.Context, models.ServiceConfig, http.Handler) error

type Middleware func(w http.ResponseWriter, r *http.Request)
type MiddlewareFactory func(Middleware) Middleware
