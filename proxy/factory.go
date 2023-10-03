package proxy

import (
	"errors"

	"github.com/hosseintrz/gaterun/config"
)

type Factory interface {
	New(cfg *config.EndpointConfig) (Proxy, error)
}

type factory struct {
	backendProxyFactory BackendProxyFactory
}

func (f factory) New(cfg *config.EndpointConfig) (proxy Proxy, err error) {
	if len(cfg.Backends) == 0 {
		err = ErrEmptyBackends
		return
	}

	if len(cfg.Backends) > 1 {
		err = errors.New("multi backends not implemented")
	}

	backend := cfg.Backends[0]
	proxy = f.backendProxyFactory(backend)
	proxy = RequestBuilderMiddleware(backend)(proxy)
	return
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
	return NewFactory(func(_ *config.BackendConfig) Proxy {
		return NoopProxy
	})
}
