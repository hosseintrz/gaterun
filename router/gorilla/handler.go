package gorilla

import (
	"context"
	"errors"
	"net/http"

	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/proxy"
)

type HandlerFactory func(*config.EndpointConfig, proxy.Proxy) http.HandlerFunc

func NewHandlerFactory(rb proxy.RequestBuilder) HandlerFactory {
	return func(ec *config.EndpointConfig, prxy proxy.Proxy) http.HandlerFunc {
		method := ec.Method
		responseWrapper := NewResponseWrapper(ec.OutputEncoding)

		return func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("GATERUN", "1.0")
			if req.Method != method {
				http.Error(rw, "methods doesn't match", http.StatusMethodNotAllowed)
				return
			}

			reqCtx, cancel := context.WithTimeout(req.Context(), ec.Timeout)
			defer cancel()

			proxyReq := rb(req, ec.QueryString, ec.TargetHeaders)
			resp, err := prxy(reqCtx, proxyReq)

			select {
			case <-reqCtx.Done():
				err = errors.New("request timeout exceeded")
			default:
			}

			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			for key, values := range resp.Metadata.Headers {
				for _, val := range values {
					rw.Header().Add(key, val)
				}
			}

			responseWrapper(rw, resp)
		}
	}
}
