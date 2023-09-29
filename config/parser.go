package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Parse(configPath string) (cfg ServiceConfig, err error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err = v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	if err = v.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		return
	}

	return
}
