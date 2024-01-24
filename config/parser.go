package config

import (
	"fmt"
	"reflect"
	"time"

	cfgmodels "github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/api/admin/models"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Parse(configPath string) (cfg cfgmodels.ServiceConfig, err error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err = v.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	hookFunc := func() mapstructure.DecodeHookFunc {
		return func(
			f reflect.Type,
			t reflect.Type,
			data interface{}) (interface{}, error) {
			if f.Kind() != reflect.String {
				return data, nil
			}
			if t != reflect.TypeOf(time.Duration(5)) {
				return data, nil
			}

			// Convert it by parsing
			duration, err := time.ParseDuration(data.(string))
			if err != nil {
				return nil, err
			}

			return models.CustomDuration(duration), nil
		}
	}

	if err = v.Unmarshal(&cfg, viper.DecodeHook(hookFunc())); err != nil {
		fmt.Printf("Error unmarshaling config: %s\n", err)
		return
	}

	return
}

// func ReadConfig(configPath string) (v *viper.Viper, err error) {
// 	v = viper.New()
// 	v.SetConfigFile(configPath)
// 	if err = viper.ReadInConfig(); err != nil {
// 		return
// 	}

// 	return
// }
