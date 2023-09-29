package proxy

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/hosseintrz/gaterun/transport/http/server"
)

type Request struct {
	Method  string
	URL     *url.URL
	Query   url.Values
	Path    string
	Body    io.ReadCloser
	Params  map[string]string
	Headers map[string][]string
}

func (r *Request) Clone() *Request {
	if r == nil {
		return nil
	}
	r2 := new(Request)
	*r2 = *r
	r2.URL = CloneURL(r.URL)
	r2.Query = cloneURLValues(r.Query)
	r2.Params = cloneParams(r.Params)
	r2.Headers = cloneHeaders(r.Headers)

	if r.Body == nil {
		return r2
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	r.Body.Close()

	r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	r2.Body = io.NopCloser(buf)

	return r2
}

func cloneParams(m map[string]string) map[string]string {
	m2 := make(map[string]string)
	for k, v := range m {
		m2[k] = v
	}
	return m2
}

func cloneHeaders(h map[string][]string) map[string][]string {
	h2 := make(map[string][]string)

	nv := 0
	for _, vv := range h {
		nv += len(vv)
	}

	tmp := make([]string, nv)
	for k, vs := range h {
		if vs == nil {
			h2[k] = nil
			continue
		}

		n := copy(tmp, vs)
		h2[k] = tmp[:n:n]
		tmp = tmp[n:]
	}

	return h2
}

func cloneURLValues(v url.Values) url.Values {
	if v == nil {
		return nil
	}
	return url.Values(http.Header(v).Clone())
}

func CloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}

	u2 := new(url.URL)
	*u2 = *u
	if u.User != nil {
		ui2 := new(url.Userinfo)
		*ui2 = *u.User
	}

	return u2
}

type RequestBuilder func(req *http.Request, queryString, headersToSend []string) *Request

type ParamExtractor func(r *http.Request) map[string]string

func NewRequestBuilder(pe ParamExtractor) RequestBuilder {
	return func(req *http.Request, queryString, headersToSend []string) *Request {
		params := pe(req)
		headers := make(map[string][]string, len(headersToSend)+3)

		for _, key := range headersToSend {
			if key == "*" {
				headers = req.Header
				break
			}

			if val, ok := req.Header[key]; ok {
				headers[key] = val
			}
		}

		headers["X-Forwarded-For"] = []string{originIP(req)}
		headers["X-Forwarded-Host"] = []string{req.Host}

		headers["Via"] = server.UserAgentHeader

		queryParams := req.URL.Query()
		query := make(map[string][]string, len(queryString))
		for _, key := range queryString {
			if key == "*" {
				query = queryParams
				break
			}

			if val, ok := queryParams[key]; ok && len(val) > 0 {
				query[key] = val
			}
		}

		return &Request{
			Method: req.Method,
			// URL:     req.URL,
			Query:   query,
			Path:    req.URL.Path,
			Body:    req.Body,
			Params:  params,
			Headers: headers,
		}
	}
}

func originIP(req *http.Request) string {
	ip := req.Header.Get("X-Forwared-For")
	ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	if len(ip) > 0 {
		return ip
	}

	if host, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return host
	}

	return ""
}
