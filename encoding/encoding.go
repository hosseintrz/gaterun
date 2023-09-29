package encoding

import (
	"encoding/json"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *map[string]interface{}) error
}

type DecoderFactory func() Decoder

func NewDecoderFactory(name string) DecoderFactory {
	switch name {
	case "no-op":
		return NoOpDecoderFactory
	default:
		return JSONDecoderFactory
	}
}

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

func NoOpDecoderFactory() Decoder {
	return &NoOpDecoder{}
}

type NoOpDecoder struct{}

func (d *NoOpDecoder) Decode(r io.Reader, v *map[string]any) error {
	return nil
}
