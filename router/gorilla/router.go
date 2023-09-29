package gorilla

import (
	"context"
	"net/http"

	"log/slog"

	gorilla "github.com/gorilla/mux"
	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/proxy"
	"github.com/hosseintrz/gaterun/router"
	"github.com/hosseintrz/gaterun/transport/http/server"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// type MiddlewareFunc func(http.Handler) http.Handler

type Config struct {
	Router         *gorilla.Router
	ProxyFactory   proxy.Factory
	HandlerFactory HandlerFactory
	Middlewares    []gorilla.MiddlewareFunc
	ServerRunner   router.ServerRunnerFunc
}

type gorillaRouter struct {
	cfg          Config
	ctx          context.Context
	serverRunner router.ServerRunnerFunc
}

func (r gorillaRouter) Run(cfg config.ServiceConfig) {
	r.cfg.Router.Use(r.cfg.Middlewares...)

	// if r.cfg.HealthCheck{
	// 	r.cfg.Router.Get("/__health", healthCheckHandler)
	// }

	r.registerEndpoints(cfg.Endpoints)

	if err := r.serverRunner(r.ctx, cfg, r.cfg.Router); err != nil {
		slog.Error(err.Error())
	}

	slog.Info("server stopped!!")
}

func (r gorillaRouter) registerEndpoints(endpoints []*config.EndpointConfig) {
	for _, endpoint := range endpoints {
		proxy, err := r.cfg.ProxyFactory.New(endpoint)
		if err != nil {
			slog.Error("error getting endpoint proxy", "endpoint", endpoint.Endpoint, "error", err)
			continue
		}
		r.registerEndpoint(endpoint, r.cfg.HandlerFactory(endpoint, proxy))
	}
}

func (r gorillaRouter) registerEndpoint(endpoint *config.EndpointConfig, handler http.HandlerFunc) {
	path := endpoint.Endpoint
	method := endpoint.Method

	r.cfg.Router.HandleFunc(path, handler).Methods(method)

	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		r.cfg.Router.HandleFunc(path, handler).Methods(method)
	default:
		slog.Error("unsupported http method: %s", method)
		return
	}

	slog.Info("registering endpoint", "method", method, "path", path)
}

type factory struct {
	cfg Config
}

func (f factory) New() router.Router {
	return f.NewWithContext(context.Background())
}

func (f factory) NewWithContext(ctx context.Context) router.Router {
	return gorillaRouter{
		cfg:          f.cfg,
		ctx:          ctx,
		serverRunner: f.cfg.ServerRunner,
	}
}

func DefaultFactory(proxyFactory proxy.Factory) factory {
	return factory{
		cfg: Config{
			Router:         gorilla.NewRouter(),
			ProxyFactory:   proxyFactory,
			Middlewares:    []gorilla.MiddlewareFunc{},
			HandlerFactory: NewHandlerFactory(proxy.NewRequestBuilder(gorillaParamsExtractor)),
			ServerRunner:   server.RunServer,
		},
	}
}

func gorillaParamsExtractor(r *http.Request) map[string]string {
	params := map[string]string{}
	title := cases.Title(language.Und)
	for key, value := range gorilla.Vars(r) {
		params[title.String(key)] = value
	}
	return params
}
