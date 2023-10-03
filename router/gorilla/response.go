package gorilla

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/hosseintrz/gaterun/proxy"
)

const (
	ApplicationJSON = "application/json"
	PlainText       = "text/plain"
)

type ResponseWrapper func(rw http.ResponseWriter, res *proxy.Response)

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

func JsonWrapper(rw http.ResponseWriter, res *proxy.Response) {
	status := res.Metadata.StatusCode
	rw.WriteHeader(status)

	//rw.Header().Set("Content-Type", "application/json")
	if res == nil {
		rw.Write([]byte{})
		return
	}

	jsonData, err := json.Marshal(res.Data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Write(jsonData)
}
func StringWrapper(rw http.ResponseWriter, res *proxy.Response) {
	status := res.Metadata.StatusCode
	rw.WriteHeader(status)

	//rw.Header().Set("Content-Type", "text/plain")
	if res == nil {
		rw.Write([]byte{})
		return
	}

	content, ok := res.Data["content"]
	if !ok {
		rw.Write([]byte{})
		return
	}

	strContent, ok := content.(string)
	if !ok {
		rw.Write([]byte{})
		return
	}

	rw.Write([]byte(strContent))
}

func NoopWrapper(rw http.ResponseWriter, res *proxy.Response) {
	if res == nil {
		rw.Write([]byte{})
		return
	}

	for key, values := range res.Metadata.Headers {
		for _, val := range values {
			rw.Header().Add(key, val)
		}
	}

	if res.Metadata.StatusCode != 0 {
		rw.WriteHeader(res.Metadata.StatusCode)
	}

	if res.Io == nil {
		return
	}

	io.Copy(rw, res.Io)
}
