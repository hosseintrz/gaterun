package proxy

import (
	"context"
	"errors"

	"github.com/hosseintrz/gaterun/config"
)

type Proxy func(context.Context, *Request) (*Response, error)

type BackendProxyFactory func(cfg *config.BackendConfig) Proxy

func NoopProxy(_ context.Context, _ *Request) (*Response, error) {
	return nil, nil
}

var (
	ErrEmptyBackends = errors.New("each endpoint should have at least one backend")
)
