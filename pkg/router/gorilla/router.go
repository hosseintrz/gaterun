package gorilla

import (
	"context"
	"fmt"
	"net/http"

	"log/slog"

	gorilla "github.com/gorilla/mux"
	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/proxy"
	"github.com/hosseintrz/gaterun/pkg/router"
	"github.com/hosseintrz/gaterun/pkg/transport/http/server"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// type MiddlewareFunc func(http.Handler) http.Handler

// func init() {

// 	admin.RegisterEndpointCallback(restartRouter)
// }

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
	cancelCtx    context.CancelFunc
	serverRunner router.ServerRunnerFunc
}

func (r *gorillaRouter) Run(cfg models.ServiceConfig) {
	r.cfg.Router.Use(r.cfg.Middlewares...)

	if cfg.HealthCheck {
		r.cfg.Router.HandleFunc("/__health", healthCheckHandler).Methods(http.MethodGet)
	}

	log.Infof("called registerEndpoints: %v", cfg.Endpoints)
	r.registerEndpoints(cfg.Endpoints)

	err := r.serverRunner(r.ctx, cfg, r.cfg.Router)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info("server stopped!!")
}

func (r *gorillaRouter) Shutdown() {
	r.cancelCtx() // Cancel the existing context
}

func (r gorillaRouter) registerEndpoints(endpoints []*models.EndpointConfig) {
	for _, endpoint := range endpoints {
		proxy, err := r.cfg.ProxyFactory.New(endpoint)
		if err != nil {
			slog.Error("error getting endpoint proxy", "endpoint", endpoint.Endpoint, "error", err)
			continue
		}
		r.registerEndpoint(endpoint, r.cfg.HandlerFactory(endpoint, proxy))
	}
}

func (r gorillaRouter) registerEndpoint(endpoint *models.EndpointConfig, handler http.HandlerFunc) {
	path := endpoint.Endpoint
	method := endpoint.Method

	//r.cfg.Router.HandleFunc(path, handler).Methods(method)

	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		r.cfg.Router.HandleFunc(path, handler).Methods(method)
	default:
		slog.Error(fmt.Sprintf("unsupported http method: %v", method))
		return
	}

	log.Infof("registering endpoint {method=%s - path=%s}", method, path)
}

// func restartRouter(r *gorillaRouter, cfg models.ServiceConfig) {
// 	log.Infoln("restrating router...")
// 	r.cancelCtx()
// 	r.Run(cfg)
// }

type factory struct {
	cfg Config
}

func (f factory) New() router.Router {
	return f.NewWithContext(context.Background())
}

func (f factory) NewWithContext(ctx context.Context) router.Router {
	ctx, cancel := context.WithCancel(ctx)
	return &gorillaRouter{
		cfg:          f.cfg,
		ctx:          ctx,
		cancelCtx:    cancel,
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

func (f *factory) AddMiddlewares(mds []router.MiddlewareFactory) {
	for _, gmd := range mds {
		md := gorillaMiddleware(gmd)
		f.cfg.Middlewares = append(f.cfg.Middlewares, md)
	}
}

func gorillaMiddleware(gmd router.MiddlewareFactory) gorilla.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(gmd(next.ServeHTTP))
	}
}

func gorillaParamsExtractor(r *http.Request) map[string]string {
	params := map[string]string{}
	title := cases.Title(language.Und)
	queryParams := r.URL.Query()
	for key, vals := range queryParams {
		params[title.String(key)] = vals[0]
	}

	// for key, value := range gorilla.Vars(r) {
	// 	params[title.String(key)] = value
	// }
	return params
}
