package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

type HTTPStatusHandler func(context.Context, *http.Response) (*http.Response, error)

func DefaultHTTPStatusHandler(ctx context.Context, res *http.Response) (*http.Response, error) {
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, ErrInvalidStatusCode
	}

	return res, nil
}

func NoOpHTTPStatusHandler(ctx context.Context, res *http.Response) (*http.Response, error) {
	return res, nil
}

func ErrorHTTPStatusHandler(ctx context.Context, resp *http.Response) (*http.Response, error) {
	if _, err := DefaultHTTPStatusHandler(ctx, resp); err == nil {
		return resp, nil
	}
	return resp, newHTTPResponseError(resp)
}

func newHTTPResponseError(resp *http.Response) HTTPResponseError {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		body = []byte{}
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	return HTTPResponseError{
		Code: resp.StatusCode,
		Msg:  string(body),
	}
}

type HTTPResponseError struct {
	Code int    `json:"http_status_code"`
	Msg  string `json:"http_body,omitempty"`
}

func (r HTTPResponseError) Error() string {
	return r.Msg
}
