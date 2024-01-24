package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Consumer struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CustomID  string    `json:"custom_id"`
	CreatedAt time.Time `json:"created_at"`
}

type EndpointRequestDTO struct {
	ID             int64
	Endpoint       string         `json:"endpoint"`
	Method         string         `json:"method"`
	Backends       StrArray       `json:"backends" gorm:"type:text[]`
	Timeout        CustomDuration `json:"timeout" gorm:"type:int"`
	CacheTTL       CustomDuration `json:"cache_ttl" gorm:"type:int"`
	QueryStrings   StrArray       `mapstructure:"query_strings" json:"query_strings" gorm:"type:text[]"`
	TargetHeaders  StrArray       `mapstructure:"target_headers" json:"target_headers" gorm:"type:text[]"`
	OutputEncoding string         `json:"output_encoding"`
}

type CustomDuration time.Duration

func (d CustomDuration) Value() (driver.Value, error) {
	val := time.Duration(d)
	return int64(val), nil
}

func (d *CustomDuration) Scan(v interface{}) error {
	if val, ok := v.(int64); ok {
		*d = CustomDuration(val)
		return nil
	}
	return fmt.Errorf("can't convert %v to customDuration", v)
}

type StrArray []string

func (a StrArray) Value() (driver.Value, error) {
	if a == nil {
		return "", nil
	}
	if len(a) == 0 {
		return "", nil
	}
	return strings.Join(a, ","), nil
}

func (a *StrArray) Scan(v interface{}) error {
	if v == nil {
		*a = []string{}
		return nil
	}

	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("error converting %v to string", v)
	}

	*a = strings.Split(str, ",")
	return nil
}

// func (cd *CustomDuration) Unmarshal(data []byte) error {
// 	var durationString string
// 	if err := json.Unmarshal(data, &durationString); err != nil {
// 		return err
// 	}

// 	duration, err := time.ParseDuration(durationString)
// 	if err != nil {
// 		return err
// 	}

// 	*cd = CustomDuration(duration)
// 	return nil
// }

// type Endpoint struct {
// 	ID             int64
// 	Endpoint       string        `mapstructure:"endpoint"`
// 	Method         string        `mapstructure:"method"`
// 	Backends       []*Backend    `mapstructure:"backends"`
// 	Timeout        time.Duration `mapstructure:"timeout"`
// 	CacheTTL       time.Duration `mapstructure:"cache_ttl"`
// 	QueryString    []string      `mapstructure:"input_query_strings"`
// 	TargetHeaders  []string      `mapstructure:"target_headers"`
// 	OutputEncoding string        `mapstructure:"output_encoding"`
// }

// type Backend struct {
// 	ID         int64
// 	Host       string   `mapstructure:"host"`
// 	Method     string   `mapstructure:"method"`
// 	URLPattern string   `mapstructure:"url_pattern"`
// 	AllowList  []string `mapstructure:"allow"`
// 	// DenyList is a set of response fields to remove. If empty, the filter id not used
// 	DenyList []string `mapstructure:"deny"`
// 	// map of response fields to be renamed and their new names
// 	Mapping map[string]string `mapstructure:"mapping"`
// 	// the encoding format
// 	Encoding string `mapstructure:"encoding"`
// 	Timeout  time.Duration
// 	// decoder to use in order to parse the received response from the API
// 	DecoderFactory string `json:"-"`
// 	// HeadersToPass defines the list of headers to pass to this backend
// 	HeadersToPass []string `mapstructure:"input_headers"`
// 	//number of concurrent calls to send to this backend
// 	ConcurrentCalls int32 `mapstructure:"concurrent_calls"`
// }

type EndpointBackend struct {
	EndpointId int64
	BackendId  int64
}
