package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceConfig ServiceConfig  `mapstructure:"service_config"`
	Database      DatabaseConfig `mapstructure:"database"`
	Logging       LoggingConfig  `mapstructure:"logging"`
}

type LoggingConfig struct {
	Level     log.Level     `mapstructure:"level"`
	Formatter log.Formatter `mapstructure:"formatter"`
	Output    io.Writer     `mapstructure:"output"`
}

type DBType string

const (
	POSTGRES DBType = "postgres"
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
	Version   int               `mapstructure:"version"`
	Name      string            `mapstructure:"name"`
	Endpoints []*EndpointConfig `mapstructure:"endpoints"`
	Host      string            `mapstructure:"host"`
	Port      int               `mapstructure:"port"`
	Timeout   time.Duration     `mapstructure:"timeout"`
	Router    RouterType        `mapstructure:"router"`
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

type RouterType string

const (
	GORILLA RouterType = "gorilla"
	GIN     RouterType = "gin"
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

func Load(configFile string) (*Config, error) {
	viper.AutomaticEnv()

	log.Infof("configFile is %s\n", configFile)

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		abPath, _ := filepath.Abs("gaterun.conf.yml")
		viper.SetConfigFile(abPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
