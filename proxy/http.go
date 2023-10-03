package proxy

import (
	"context"
	"net/http"
	"strconv"

	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/encoding"
	"github.com/hosseintrz/gaterun/transport/http/client"
)

var httpProxy = CustomHTTPProxyFactory(client.DefaultHttpClientFactory{})

func DefaultHTTPProxyFactory(c *http.Client) BackendProxyFactory {
	return CustomHTTPProxyFactory(client.DefaultHttpClientFactory{})
}

func CustomHTTPProxyFactory(cf client.HTTPClientFactory) BackendProxyFactory {
	return func(cfg *config.BackendConfig) Proxy {
		return NewHTTPProxy(cfg, cf)
	}
}

func NewHTTPProxy(cfg *config.BackendConfig, cf client.HTTPClientFactory) Proxy {
	// return NewHTTPProxyWithHTTPExecutor(cfg, client.DefaultHTTPReqeustExecutor(cf), encoding.NewDecoderFactory(cfg.DecoderFactory))
	return NewHTTPProxyWithHTTPExecutor(cfg, client.DefaultHTTPReqeustExecutor(cf), encoding.CustomDecoderFactory)

}

func NewHTTPProxyWithHTTPExecutor(cfg *config.BackendConfig, requestExecutor client.HTTPRequestExecutor, decoderFactory encoding.DecoderFactory) Proxy {
	responseParser := NewDefaultHTTPResponseParser(decoderFactory)
	statusHandler := client.DefaultHTTPStatusHandler
	return NewCustomHTTPProxy(cfg, requestExecutor, statusHandler, responseParser)
}

func NewCustomHTTPProxy(
	cfg *config.BackendConfig,
	reqExecutor client.HTTPRequestExecutor,
	statusHandler client.HTTPStatusHandler,
	responseParser HTTPResponseParser) Proxy {
	return func(ctx context.Context, req *Request) (*Response, error) {
		request, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
		if err != nil {
			return nil, err
		}

		headers := make(map[string][]string, len(req.Headers))
		for key, values := range headers {
			newHeaders := make([]string, len(values))
			copy(newHeaders, values)
			request.Header[key] = values
		}

		if req.Body != nil {
			if v, ok := req.Headers["Content-Length"]; ok && len(v) == 1 && v[0] != "chunked" {
				if size, err := strconv.Atoi(v[0]); err == nil {
					request.ContentLength = int64(size)
				}
			}
		}

		response, err := reqExecutor.Execute(ctx, request)
		if err != nil {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		if err != nil {
			return nil, err
		}

		// resp, err := statusHandler(ctx, response)
		// if err != nil {
		// 	return nil, err
		// }

		return responseParser(ctx, response)
	}
}
