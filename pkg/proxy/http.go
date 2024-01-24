package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/encoding"
	"github.com/hosseintrz/gaterun/pkg/transport/http/client"
	"github.com/sirupsen/logrus"
)

var httpProxy = CustomHTTPProxyFactory(client.DefaultHttpClientFactory{})

func DefaultHTTPProxyFactory(c *http.Client) BackendProxyFactory {
	return CustomHTTPProxyFactory(client.DefaultHttpClientFactory{})
}

func CustomHTTPProxyFactory(cf client.HTTPClientFactory) BackendProxyFactory {
	return func(cfg *models.BackendConfig) Proxy {
		return NewHTTPProxy(cfg, cf)
	}
}

func NewHTTPProxy(cfg *models.BackendConfig, cf client.HTTPClientFactory) Proxy {
	// return NewHTTPProxyWithHTTPExecutor(cfg, client.DefaultHTTPReqeustExecutor(cf), encoding.NewDecoderFactory(cfg.DecoderFactory))
	return NewHTTPProxyWithHTTPExecutor(cfg, client.DefaultHTTPReqeustExecutor(cf), encoding.CustomDecoderFactory)

}

func NewHTTPProxyWithHTTPExecutor(cfg *models.BackendConfig, requestExecutor client.HTTPRequestExecutor, decoderFactory encoding.DecoderFactory) Proxy {
	responseParser := NewDefaultHTTPResponseParser(decoderFactory)
	statusHandler := client.DefaultHTTPStatusHandler
	return NewCustomHTTPProxy(cfg, requestExecutor, statusHandler, responseParser)
}

func NewCustomHTTPProxy(
	cfg *models.BackendConfig,
	reqExecutor client.HTTPRequestExecutor,
	statusHandler client.HTTPStatusHandler,
	responseParser HTTPResponseParser) Proxy {
	return func(ctx context.Context, req *Request) (*Response, error) {
		url := createURL(req.URL)
		request, err := http.NewRequest(req.Method, url, req.Body)
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

		// q := req.Query
		// for k, values := range q {
		// 	q.Add(k, strings.Join(values, ","))
		// }
		request.URL.RawQuery = req.Query.Encode()

		logrus.Infof("emit %s request to %s", req.Method, url)
		response, err := reqExecutor.Execute(ctx, request)
		if err != nil {
			return nil, err
		}
		logrus.Infof("received response %v", response)

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

func createURL(url *url.URL) string {
	return fmt.Sprintf("%s%s", url.Host, url.Path)
}
