package gateway

import (
	"fmt"

	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/api/admin"
	"github.com/hosseintrz/gaterun/pkg/auth"
	"github.com/hosseintrz/gaterun/pkg/proxy"
	"github.com/hosseintrz/gaterun/pkg/ratelimit"
	"github.com/hosseintrz/gaterun/pkg/router"
	"github.com/hosseintrz/gaterun/pkg/router/gorilla"

	log "github.com/sirupsen/logrus"
)

type Gateway struct {
	Conf   models.ServiceConfig
	Router router.Router
}

func NewGateway(conf models.ServiceConfig) *Gateway {
	return &Gateway{
		Conf: conf,
	}
}

func (g *Gateway) Start() {
	routerFactory := createRouterFactory(g.Conf)

	g.Router = routerFactory.New()

	admin.RegisterEndpointCallback(func(conf models.ServiceConfig) error {
		g.Router.Shutdown()

		log.Infoln("after shutting down in gatway")

		g.Conf = conf
		config.SetGlobalServiceConf(conf)

		go g.Start()

		log.Infoln("after starting the gateway again")

		return nil
	})

	g.Router.Run(g.Conf)
	log.Infoln("returning from gateway.Start()")
}

func createRouterFactory(cfg models.ServiceConfig) router.Factory {
	middlewares := assembleMiddlewares(cfg)

	switch cfg.Router {
	case models.GORILLA:
		r := gorilla.DefaultFactory(proxy.NewDefaultFactory())
		r.AddMiddlewares(middlewares)
		return r
	case models.GIN:
		log.Error("not implemented")
	default:
		log.Error("invalid router name")
	}

	return nil
}

func assembleMiddlewares(cfg models.ServiceConfig) []router.MiddlewareFactory {
	mds := make([]router.MiddlewareFactory, 0)

	authMD := getAuthMD(cfg.AuthType)
	if authMD != nil {
		mds = append(mds, authMD)
	}

	rateLimitMD := getRateLimiterMD(*cfg.RateLimit)
	if rateLimitMD != nil {
		mds = append(mds, rateLimitMD)
	}

	return mds
}

func getAuthMD(authType models.AuthType) router.MiddlewareFactory {
	switch authType {
	case models.API_KEY:
		return auth.APIKeyMiddleware
	case models.BASIC:
		return auth.BasicMiddleware
	case models.JWT:
		return auth.JWTMiddleware
	}

	return nil
}

func getRateLimiterMD(cfg models.RateLimitConfig) router.MiddlewareFactory {
	switch cfg.Algorithm {
	case models.TokenBucket:
		return ratelimit.TokenBucketMiddlewareWithCfg(cfg)
	case models.FixedWindowCounter:
		fmt.Println("this rate limit algorithm is not implemented")
		return nil
	case models.SlidingWindowLog:
		return ratelimit.SlidingWindowLogMDFactory(cfg)
		//fmt.Println("this rate limit algorithm is not implemented")
		//return nil
	case models.SlidingWindowCounter:
		fmt.Println("this rate limit algorithm is not implemented")
		return nil
	}

	return nil
}
