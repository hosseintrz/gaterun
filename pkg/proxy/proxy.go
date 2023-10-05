package proxy

import (
	"context"
	"errors"
	"net/http"

	"github.com/hosseintrz/gaterun/config"
)

type Proxy func(context.Context, *Request) (*Response, error)

type BackendProxyFactory func(cfg *config.BackendConfig) Proxy

func NoopProxy(_ context.Context, _ *Request) (*Response, error) {
	return &Response{
		IsComplete: true,
		Data: map[string]interface{}{
			"ok": true,
		},
		Metadata: Metadata{
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
			StatusCode: http.StatusOK,
		},
	}, nil
}

var (
	ErrEmptyBackends = errors.New("each endpoint should have at least one backend")
)
