package gorilla

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/proxy"
)

type HandlerFactory func(*models.EndpointConfig, proxy.Proxy) http.HandlerFunc

func NewHandlerFactory(rb proxy.RequestBuilder) HandlerFactory {
	return func(ec *models.EndpointConfig, prxy proxy.Proxy) http.HandlerFunc {
		method := ec.Method

		return func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Set("GATERUN", "1.0")
			if req.Method != method {
				http.Error(rw, "methods doesn't match", http.StatusMethodNotAllowed)
				return
			}

			reqCtx, cancel := context.WithTimeout(req.Context(), time.Duration(ec.Timeout))
			defer cancel()

			urlParams, err := matchParams(ec.Endpoint, req.URL.Host+req.URL.Path)
			if err != nil {
				return
			}
			proxyReq := rb(req, ec.QueryStrings, ec.TargetHeaders, urlParams)
			// logrus.Infof("sending request to %s%s", proxyReq.URL.Host, proxyReq.URL.Path)
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

			responseWrapper := NewResponseWrapper(resp.Metadata.Headers, ec.OutputEncoding)

			response, err := responseWrapper(&rw, resp)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			// for key, values := range resp.Metadata.Headers {
			// 	for _, val := range values {
			// 		rw.Header().Add(key, val)
			// 	}
			// }

			rw.WriteHeader(resp.Metadata.StatusCode)
			_, err = rw.Write(response)
			if err != nil {
				return
			}
		}
	}
}

func matchParams(endpoint, reqURL string) (map[string]string, error) {
	// Create a regular expression for matching {param} in the endpoint
	re := regexp.MustCompile(`{([^}]+)}`)

	// Find all matches in the endpoint
	matches := re.FindAllStringSubmatch(endpoint, -1)

	// Extract parameter names from matches
	paramNames := make([]string, len(matches))
	for i, match := range matches {
		if len(match) > 1 {
			paramNames[i] = match[1]
		}
	}

	// Create a regular expression pattern for matching the endpoint
	// Replace {param} with a capturing group for the parameter value
	//pattern := regexp.QuoteMeta(endpoint)
	pattern := re.ReplaceAllString(endpoint, `([^/]+)`)
	re = regexp.MustCompile(pattern)

	// Match the request URL against the pattern
	matches2 := re.FindStringSubmatch(reqURL)

	// Extract parameter values and create a map
	params := make(map[string]string)
	for i, name := range paramNames {
		if i+1 < len(matches2) {
			params[name] = matches2[i+1]
		}
	}

	return params, nil
}

func healthCheckHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(`{"status":"ok"}`))
}
