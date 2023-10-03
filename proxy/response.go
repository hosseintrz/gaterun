package proxy

import (
	"context"
	"io"
	"net/http"

	"github.com/hosseintrz/gaterun/encoding"
)

type Metadata struct {
	Headers    map[string][]string
	StatusCode int
}

type Response struct {
	Data       map[string]interface{}
	IsComplete bool
	Metadata   Metadata
	Io         io.Reader
}

type HTTPResponseParser func(context.Context, *http.Response) (*Response, error)

func NewDefaultHTTPResponseParser(decoderFactory encoding.DecoderFactory) HTTPResponseParser {
	return func(ctx context.Context, res *http.Response) (*Response, error) {
		defer res.Body.Close()

		data := make(map[string]interface{})
		decoder := decoderFactory(res.Header)
		// decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(res.Body, &data); err != nil {
			return nil, err
		}

		response := &Response{Data: data, IsComplete: true}

		response.Metadata.Headers = cloneHeaders(res.Header)
		response.Metadata.StatusCode = res.StatusCode

		return response, nil
	}
}
