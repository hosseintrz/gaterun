package gorilla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hosseintrz/gaterun/pkg/proxy"
)

const (
	ApplicationJSON = "application/json"
	PlainText       = "text/plain"
)

type ResponseWrapper func(rw *http.ResponseWriter, res *proxy.Response) ([]byte, error)

func NewResponseWrapper(headers map[string][]string, outputEncoding string) ResponseWrapper {
	encoding := outputEncoding
	contentType := getContentType(headers)
	if len(contentType) > 0 {
		encoding = contentType
	}

	switch encoding {
	case ApplicationJSON:
		return JsonWrapper
	case PlainText:
		return StringWrapper
	default:
		return NoopWrapper
	}
}

func getContentType(headers map[string][]string) string {
	// Check if the Content-Type header is present
	contentTypeValues, ok := headers["Content-Type"]
	if !ok || len(contentTypeValues) == 0 {
		return "" // Content-Type not found
	}

	// Content-Type header found, extract the first value
	contentType := contentTypeValues[0]

	// Extract only the mime type, removing any additional parameters (e.g., charset)
	contentType = strings.Split(contentType, ";")[0]

	return strings.TrimSpace(contentType)
}

func JsonWrapper(rw *http.ResponseWriter, res *proxy.Response) (response []byte, err error) {
	//rw.Header().Set("Content-Type", "application/json")
	if res == nil {
		return
	}

	response, err = json.Marshal(res.Data)
	if err != nil {
		return
	}

	return response, nil
}

func StringWrapper(rw *http.ResponseWriter, res *proxy.Response) (response []byte, err error) {
	//rw.Header().Set("Content-Type", "text/plain")
	if res == nil {
		err = fmt.Errorf("response is nil")
		return
	}

	content, ok := res.Data["content"]
	if !ok {
		err = fmt.Errorf("content in response doesn't exists")
		return
	}

	strContent, ok := content.(string)
	if !ok {
		err = fmt.Errorf("invalid repsonse type, not string")
		return
	}

	response = []byte(strContent)
	return
}

func NoopWrapper(rw *http.ResponseWriter, res *proxy.Response) (response []byte, err error) {
	if res == nil {
		err = fmt.Errorf("response is nil")
		return
	}

	if res.Io == nil {
		err = fmt.Errorf("IO is nil")
		return
	}

	io.Copy((*rw), res.Io)
	return nil, nil
}
