package gateway

import (
	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/pkg/proxy"
	"github.com/hosseintrz/gaterun/pkg/router"
	"github.com/hosseintrz/gaterun/pkg/router/gorilla"
	log "github.com/sirupsen/logrus"
)

type Gateway struct {
	Conf config.ServiceConfig
}

func NewGateway(conf config.ServiceConfig) *Gateway {
	return &Gateway{
		Conf: conf,
	}
}

func (g *Gateway) Start() {
	routerFactory := getRouterFactory(g.Conf.Router)
	routerFactory.New().Run(g.Conf)
}

func getRouterFactory(routerName config.RouterType) router.Factory {
	switch routerName {
	case config.GORILLA:
		return gorilla.DefaultFactory(proxy.NewDefaultFactory())
	case config.GIN:
		log.Error("not implemented")
	default:
		log.Error("invalid router name")
	}

	return nil
}
