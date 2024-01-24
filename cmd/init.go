package cmd

import (
	"context"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/config/persistence"
	"github.com/hosseintrz/gaterun/pkg/cache/redis"
	"github.com/hosseintrz/gaterun/pkg/database"
	log "github.com/sirupsen/logrus"
)

var (
	globalConfig *models.Config
)

func initConfig(configFile string) {
	var err error
	globalConfig, err = persistence.Load(configFile)
	if err != nil {
		log.WithError(err).Fatal("couldn't load config file")
	}
}

func initLogger() {
	logCfg := globalConfig.Logging

	log.SetLevel(logCfg.Level)
	log.SetFormatter(logCfg.Formatter)
	log.SetOutput(logCfg.Output)
}

func initDatabase() {
	dbCfg := globalConfig.Database

	err := database.InitDB(dbCfg)
	if err != nil {
		log.WithError(err).Error("couldn't connect to database")
	}
}

func initRedis(ctx context.Context) {
	cfg := globalConfig.Redis
	err := redis.InitRedis(ctx, cfg)
	if err != nil {
		log.WithError(err).Error("couldn't connect to redis")
	}
}
