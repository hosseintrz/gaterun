package gorilla

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hosseintrz/gaterun/proxy"
)

type ResponseWrapper func(rw http.ResponseWriter, res *proxy.Response)

func NewResponseWrapper(outputEncoding string) ResponseWrapper {
	switch outputEncoding {
	case "json":
		return JsonWrapper
	case "string":
		return StringWrapper
	default:
		return NoopWrapper
	}
}

func JsonWrapper(rw http.ResponseWriter, res *proxy.Response) {
	status := res.Metadata.StatusCode
	rw.WriteHeader(status)

	rw.Header().Set("Content-Type", "application/json")
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

	rw.Header().Set("Content-Type", "text/plain")
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
