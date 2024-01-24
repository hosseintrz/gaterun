package config

import "github.com/hosseintrz/gaterun/config/models"

var (
	globalConfig *models.Config
)

func SetGlobalConf(conf *models.Config) {
	globalConfig = conf
}

func SetGlobalServiceConf(conf models.ServiceConfig) {
	globalConfig.ServiceConfig = conf
}

func GetGlobalConfig() *models.Config {
	return globalConfig
}
