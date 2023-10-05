package proxy

import (
	"context"
	"net/url"

	"github.com/hosseintrz/gaterun/config"
)

type Middleware func(Proxy) Proxy

func RequestBuilderMiddleware(cfg *config.BackendConfig) Middleware {
	return func(prxy Proxy) Proxy {
		return func(ctx context.Context, req *Request) (*Response, error) {
			r2 := req.Clone()
			parsedUrl, err := url.Parse(cfg.Host)
			if err != nil {
				return nil, err
			}
			r2.URL = parsedUrl
			// r2.URL.Path = cfg.URLPattern
			// path := generatePath()
			r2.URL.Path = req.Path
			r2.Path = req.Path
			r2.Method = cfg.Method
			return prxy(ctx, r2)
		}
	}
}
