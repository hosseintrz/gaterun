package cmd

import (
	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/pkg/database"
	log "github.com/sirupsen/logrus"
)

var (
	globalConfig *config.Config
)

func initConfig(configFile string) {
	var err error
	globalConfig, err = config.Load(configFile)
	if err != nil {
		log.WithError(err).Fatal("couldn't load config file")
	}
}

func initLogger() {
	logCfg := globalConfig.Logging

	log.Info("before")

	log.SetLevel(logCfg.Level)
	log.Info("before1")

	log.SetFormatter(logCfg.Formatter)
	log.Info("before2")

	log.SetOutput(logCfg.Output)
	log.Info("after")
}

func initDatabase() {
	dbCfg := globalConfig.Database

	err := database.InitDatabase(dbCfg)
	if err != nil {
		log.WithError(err).Error("couldn't initialize database")
	}
}
