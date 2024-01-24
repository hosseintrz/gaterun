package models

import (
	"io"
	"os"
	"time"

	"github.com/hosseintrz/gaterun/pkg/api/admin/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceConfig ServiceConfig  `mapstructure:"service_config"`
	Database      DatabaseConfig `mapstructure:"database"`
	Logging       LoggingConfig  `mapstructure:"logging"`
	Redis         RedisConfig    `mapstructure:"redis"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LoggingConfig struct {
	Level     log.Level     `mapstructure:"level"`
	Formatter log.Formatter `mapstructure:"formatter"`
	Output    io.Writer     `mapstructure:"output"`
}

type DBType string

const (
	POSTGRES  DBType = "postgres"
	CASSANDRA DBType = "cassandra"
)

type DatabaseConfig struct {
	Type     DBType        `mapstructure:"type"`
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Timeout  time.Duration `mapstructure:"timeout"`
	DbName   string        `mapstructure:"db_name"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	SslMode  string        `mapstructure:"ssl_mode"`
	Schema   string        `mapstructure:"schema"`
}

type ServiceConfig struct {
	Version     int               `mapstructure:"version"`
	Name        string            `mapstructure:"name"`
	Endpoints   []*EndpointConfig `mapstructure:"endpoints"`
	Host        string            `mapstructure:"host"`
	Port        int               `mapstructure:"port"`
	Timeout     time.Duration     `mapstructure:"timeout"`
	Router      RouterType        `mapstructure:"router"`
	HealthCheck bool              `mapstructure:"health_check"`
	AuthType    AuthType          `mapstructure:"auth_type"`
	RateLimit   *RateLimitConfig  `mapstructure:"rate_limit"`
}

type EndpointConfig struct {
	ID             int64
	Endpoint       string           `mapstructure:"endpoint" json:"endpoint"`
	Method         string           `mapstructure:"method" json:"method"`
	Backends       []*BackendConfig `mapstructure:"backends" json:"backends"`
	Timeout        time.Duration    `mapstructure:"timeout" json:"timeout" gorm:"type:int"`
	CacheTTL       time.Duration    `mapstructure:"cache_ttl" json:"cache_ttl" gorm:"type:int"`
	QueryStrings   models.StrArray  `mapstructure:"input_query_strings" json:"query_string"`
	TargetHeaders  models.StrArray  `mapstructure:"target_headers" json:"target_headers"`
	OutputEncoding string           `mapstructure:"output_encoding" json:"output_encoding"`
	RateLimit      RateLimitConfig  `mapstructure:"rate_limit"`
}

type BackendConfig struct {
	ID         int64
	Host       string          `mapstructure:"host" json:"host"`
	Method     string          `mapstructure:"method" json:"method"`
	URLPattern string          `mapstructure:"url_pattern" json:"url_pattern"`
	AllowList  models.StrArray `mapstructure:"allow_list" json:"allow_list" gorm:"type:text[]"`
	// DenyList is a set of response fields to remove. If empty, the filter id not used
	DenyList models.StrArray `mapstructure:"deny" json:"deny" gorm:"type:text[]"`
	// map of response fields to be renamed and their new names
	Mapping map[string]string `mapstructure:"mapping" json:"mapping"`
	// the encoding format
	Encoding string                `mapstructure:"encoding" json:"encoding"`
	Timeout  models.CustomDuration `json:"timeout" gorm:"type:int"`
	// decoder to use in order to parse the received response from the API
	DecoderFactory string `mapstructure:"-" json:"decoder_factory" gorm:"-"`
	// HeadersToPass defines the list of headers to pass to this backend
	HeadersToPass models.StrArray `mapstructure:"headers_to_pass" json:"headers_to_pass"`
	//number of concurrent calls to send to this backend
	ConcurrentCalls int32 `mapstructure:"concurrent_calls" json:"concurrent_calls"`
	//rate limit config
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}

type RouterType string

const (
	GORILLA RouterType = "gorilla"
	GIN     RouterType = "gin"
)

type AuthType string

const (
	API_KEY AuthType = "apikey"
	BASIC   AuthType = "basic"
	JWT     AuthType = "jwt"
)

type RateLimitConfig struct {
	Domain    RLDomain      `mapstructure:"domain"`
	Interval  time.Duration `mapstructure:"interval"`
	Threshold int           `mapstrcture:"threshold"`
	Scope     RLScope       `mapstructure:"scope"`
	Target    RLTarget      `mapstructure:"target"`
	Algorithm RLAlgorithm   `mapstructure:"algorithm"`
}

type RLDomain string
type RLScope int32
type RLTarget int32
type RLAlgorithm int32

const (
	DomainService  RLDomain = "service"
	DomainEndpoint RLDomain = "endpoint"
	DomainBackend  RLDomain = "backend"
)

// const (
// 	UnitSecond RLUnit = "second"
// 	UnitMinute RLUnit = "minute"
// 	UnitHour   RLUnit = "hour"
// 	UnitDay    RLUnit = "day"
// )

const (
	UserScope RLScope = iota + 1
	GlobalScope
)

const (
	TargetIP RLTarget = iota + 1
	TargetID
)

const (
	TokenBucket RLAlgorithm = iota + 1
	FixedWindowCounter
	SlidingWindowLog
	SlidingWindowCounter
)

func init() {
	viper.SetDefault("serviceConfig.version", 3)
	viper.SetDefault("serviceConfig.name", "gaterun api-gateway")
	viper.SetDefault("serviceConfig.port", 8000)
	viper.SetDefault("serviceConfig.cache_ttl", 3000)
	viper.SetDefault("serviceConfig.timeout", 4)
	viper.SetDefault("serviceConfig.host", "localhost")
	viper.SetDefault("serviceConfig.router", GORILLA)

	viper.SetDefault("database.type", POSTGRES)
	viper.SetDefault("database.host", "127.0.0.1")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.timeout", 5000)
	viper.SetDefault("database.db_name", "gaterun")
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "pass1234")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.schema", "")

	viper.SetDefault("logging.level", log.InfoLevel)
	viper.SetDefault("logging.formatter", &log.JSONFormatter{})
	viper.SetDefault("logging.output", os.Stdout)
}
