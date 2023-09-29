package client

import (
	"context"
	"net/http"
)

var defaultHTTPClient = &http.Client{}

// type HTTPClientFactory func(context.Context) *http.Client

// func NewHTTPClient(_ context.Context) *http.Client {
// 	return defaultHTTPClient
// }

type HTTPClientFactory interface {
	New(context.Context) *http.Client
}

type DefaultHttpClientFactory struct{}

func (cf DefaultHttpClientFactory) New(ctx context.Context) *http.Client {
	return defaultHTTPClient
}

type HTTPRequestExecutor interface {
	Execute(context.Context, *http.Request) (*http.Response, error)
}

type httpRequestExecutor struct {
	clientFactory HTTPClientFactory
}

func (e httpRequestExecutor) Execute(ctx context.Context, req *http.Request) (*http.Response, error) {
	return e.clientFactory.New(ctx).Do(req.WithContext(ctx))
}

func DefaultHTTPReqeustExecutor(clientFactory HTTPClientFactory) HTTPRequestExecutor {
	return &httpRequestExecutor{
		clientFactory: clientFactory,
	}
}
