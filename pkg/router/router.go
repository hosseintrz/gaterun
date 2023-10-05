package router

import (
	"context"
	"net/http"

	"github.com/hosseintrz/gaterun/config"
)

type Router interface {
	Run(config.ServiceConfig)
}

type RouterFunc func(config.ServiceConfig)

func (f RouterFunc) Run(cfg config.ServiceConfig) { f(cfg) }

type Factory interface {
	New() Router
	NewWithContext(context.Context) Router
}

type ServerRunnerFunc func(context.Context, config.ServiceConfig, http.Handler) error
