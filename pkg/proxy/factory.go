package proxy

import (
	"errors"

	"github.com/hosseintrz/gaterun/config/models"
)

type Factory interface {
	New(cfg *models.EndpointConfig) (Proxy, error)
}

type factory struct {
	backendProxyFactory BackendProxyFactory
}

func (f factory) New(cfg *models.EndpointConfig) (proxy Proxy, err error) {
	if len(cfg.Backends) == 0 {
		err = ErrEmptyBackends
		return
	}

	if len(cfg.Backends) > 1 {
		err = errors.New("multi backends not implemented")
	}

	backend := cfg.Backends[0]
	proxy = f.backendProxyFactory(backend)
	proxy = applyMiddlewares(proxy, backend, cfg.Endpoint,
		RequestBuilderMiddleware,
	)

	if backend.ConcurrentCalls > 1 {
		proxy = applyMiddleware(proxy, backend, ConcurrentMiddleware)
	}

	// proxy = applyMiddlewares(f, backend,
	// 	RequestBuilderMiddleware,
	// 	ConcurrentMiddleware,
	// )
	// proxy = f.backendProxyFactory(backend)
	// proxy = RequestBuilderMiddleware(backend)(proxy)
	return
}

// func applyMiddlewares(f factory, backend *config.BackendConfig, mds ...func(*config.BackendConfig) Middleware) Proxy {
// 	proxy := f.backendProxyFactory(backend)
// 	for _, md := range mds {
// 		mdProxy := md(backend)
// 		proxy = mdProxy(proxy)
// 	}

// 	return proxy
// }

func applyMiddleware(proxy Proxy, backend *models.BackendConfig, middlewareFactory func(*models.BackendConfig) Middleware) Proxy {
	md := middlewareFactory(backend)
	return md(proxy)
}

func applyMiddlewares(proxy Proxy, backend *models.BackendConfig, endpoint string, mds ...MiddlewareFactory) Proxy {
	for _, md := range mds {
		mdProxy := md(backend, endpoint)
		proxy = mdProxy(proxy)
	}

	return proxy
}

func NewFactory(backendProxyFactory BackendProxyFactory) Factory {
	return factory{
		backendProxyFactory: backendProxyFactory,
	}
}

func NewDefaultFactory() Factory {
	return NewFactory(httpProxy)
}

func NoOpProxyFactory() Factory {
	return NewFactory(func(_ *models.BackendConfig) Proxy {
		return NoopProxy
	})
}
