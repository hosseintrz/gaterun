package config

import (
	"time"
)

type ServiceConfig struct {
	Version   int               `mapstructure:"version"`
	Name      string            `mapstructure:"name"`
	Endpoints []*EndpointConfig `mapstructure:"endpoints"`
	Host      string            `mapstructure:"host"`
	Port      int               `mapstructure:"port"`
	Timeout   time.Duration     `mapstructure:"timeout"`
}

type EndpointConfig struct {
	Endpoint       string           `mapstructure:"endpoint"`
	Method         string           `mapstructure:"method"`
	Backends       []*BackendConfig `mapstructure:"backends"`
	Timeout        time.Duration    `mapstructure:"timeout"`
	CacheTTL       time.Duration    `mapstructure:"cache_ttl"`
	QueryString    []string         `mapstructure:"input_query_strings"`
	TargetHeaders  []string         `mapstructure:"target_headers"`
	OutputEncoding string           `mapstructure:"output_encoding"`
}

type BackendConfig struct {
	Host       string   `mapstructure:"host"`
	Method     string   `mapstructure:"method"`
	URLPattern string   `mapstructure:"url_pattern"`
	AllowList  []string `mapstructure:"allow"`
	// DenyList is a set of response fields to remove. If empty, the filter id not used
	DenyList []string `mapstructure:"deny"`
	// map of response fields to be renamed and their new names
	Mapping map[string]string `mapstructure:"mapping"`
	// the encoding format
	Encoding string `mapstructure:"encoding"`
	Timeout  time.Duration
	// decoder to use in order to parse the received response from the API
	DecoderFactory string `json:"-"`
	// HeadersToPass defines the list of headers to pass to this backend
	HeadersToPass []string `mapstructure:"input_headers"`
}
