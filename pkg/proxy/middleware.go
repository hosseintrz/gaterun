package proxy

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/hosseintrz/gaterun/config/models"
)

type Middleware func(...Proxy) Proxy
type MiddlewareFactory func(cfg *models.BackendConfig, endpoint string) Middleware

func RequestBuilderMiddleware(cfg *models.BackendConfig, endpoint string) Middleware {
	return func(next ...Proxy) Proxy {
		return func(ctx context.Context, req *Request) (*Response, error) {
			r2 := req.Clone()
			parsedUrl, err := url.Parse(cfg.Host)
			if err != nil {
				return nil, err
			}
			r2.URL = parsedUrl
			// r2.URL.Path = cfg.URLPattern
			// path := generatePath()
			parsedUrlPattern := parseUrlParams(cfg.URLPattern, req.URLParams)
			parsedEndpoint := parseUrlParams(endpoint, req.URLParams)
			path := parsedUrlPattern + strings.TrimLeft(req.Path, parsedEndpoint)

			// path := strings.TrimLeft(req.Path, endpoint)
			// path = fmt.Sprintf("%s/%s", cfg.URLPattern, path)

			r2.URL.Host = cfg.Host
			r2.URL.Path = path
			r2.Path = path
			r2.Method = cfg.Method
			return next[0](ctx, r2)
		}
	}
}

func parseUrlParams(url string, params map[string]string) string {
	re := regexp.MustCompile(`{([^}]+)}`)

	result := re.ReplaceAllStringFunc(url, func(match string) string {
		paramName := match[1 : len(match)-1]

		if val, ok := params[paramName]; ok {
			return val
		}

		return match
	})

	return result
}
