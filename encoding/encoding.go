package encoding

import (
	"encoding/json"
	"io"
	"strings"
)

type Decoder interface {
	Decode(io.Reader, *map[string]interface{}) error
}

type DecoderFactory func(map[string][]string) Decoder

func JSONDecoderFactory() Decoder {
	return &JSONDecoder{}
}

type JSONDecoder struct{}

func (jd *JSONDecoder) Decode(r io.Reader, v *map[string]interface{}) error {
	dataBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dataBytes, v)
	if err != nil {
		return err
	}

	return nil
}

type StringDecoder struct{}

func (d *StringDecoder) Decode(r io.Reader, v *map[string]interface{}) error {
	dataBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	data := string(dataBytes)
	(*v)["content"] = data

	return nil
}

func NoOpDecoderFactory() Decoder {
	return &NoOpDecoder{}
}

type NoOpDecoder struct{}

func (d *NoOpDecoder) Decode(r io.Reader, v *map[string]any) error {
	return nil
}

func CustomDecoderFactory(headers map[string][]string) Decoder {
	contentType := getContentType(headers)
	switch contentType {
	case "application/json":
		return &JSONDecoder{}
	case "text/plain":
		return &StringDecoder{}
	default:
		return &NoOpDecoder{}
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
