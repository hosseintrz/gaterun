package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/api/admin/persistence"
	"github.com/hosseintrz/gaterun/pkg/database"
	"github.com/hosseintrz/gaterun/pkg/database/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func Load(configFile string) (*models.Config, error) {
	viper.AutomaticEnv()

	logrus.Infof("configFile is %s\n", configFile)

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		abPath, _ := filepath.Abs("gaterun.conf.yml")
		viper.SetConfigFile(abPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	var config models.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadConfigFromDb(ctx context.Context) (conf *models.ServiceConfig, err error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return
	}

	var jsonData string
	err = db.Table("configs").Select("data").Take(&jsonData).Error
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(jsonData), &conf)
	if err != nil {
		return
	}

	return conf, nil
}

func AssembleConfig(ctx context.Context) (cfg *models.ServiceConfig, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	return assembleConfigTx(ctx, db)
}

func assembleConfigTx(ctx context.Context, tx *gorm.DB) (cfg *models.ServiceConfig, err error) {
	cfg, err = LoadConfigFromDb(ctx)
	if err != nil {
		return
	}

	endpointIds := []int64{}
	err = tx.Table("endpoint").Select("id").Scan(&endpointIds).Error
	if err != nil {
		return
	}

	endpointMap := fetchToMap(cfg.Endpoints)

	for _, id := range endpointIds {
		endpoint, err := persistence.FetchEndpointTx(ctx, tx, id)
		if err != nil {
			return nil, err
		}

		if _, exists := endpointMap[endpoint.ID]; !exists {
			cfg.Endpoints = append(cfg.Endpoints, endpoint)
			endpointMap[endpoint.ID] = endpoint
		}
	}

	return
}

func fetchToMap(endpoints []*models.EndpointConfig) map[int64]*models.EndpointConfig {
	m := make(map[int64]*models.EndpointConfig)
	for _, endpoint := range endpoints {
		if _, ok := m[endpoint.ID]; !ok {
			m[endpoint.ID] = endpoint
		}
	}
	return m
}

func SaveConfigToDB(ctx context.Context, conf *models.Config) (err error) {
	db, err := postgres.GetDB(ctx)
	if err != nil {
		return
	}

	jsonData, err := json.Marshal(conf.ServiceConfig)
	if err != nil {
		return
	}

	var configId int64
	err = db.Table("configs").Select("id").Take(&configId).Error
	if err != nil {
		return
	}

	model := &struct {
		ID   int64
		Data string
	}{
		Data: string(jsonData),
	}

	if configId > 0 {
		err = db.Table("configs").Where("id", configId).Updates(model).Error
	} else {
		err = db.Table("configs").Create(model).Error
	}
	if err != nil {
		return
	}

	return nil
}
