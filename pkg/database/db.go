package database

import (
	"context"
	"fmt"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/database/cassandra"
	"github.com/hosseintrz/gaterun/pkg/database/postgres"
	"gorm.io/gorm"
)

var dbType string

// type Database interface {
// 	InitDatabase(cfg config.DatabaseConfig) error
// }

func InitDB(cfg models.DatabaseConfig) error {
	switch cfg.Type {
	case "postgres":
		postgres.InitDatabase(cfg)
	case "cassandra":
		cassandra.InitDatabase(cfg)
	default:
		return fmt.Errorf("invalid database type")
	}

	dbType = string(cfg.Type)

	return nil
}

func GetDB(ctx context.Context) (*gorm.DB, error) {
	switch dbType {
	case "postgres":
		return postgres.GetDB(ctx)
	default:
		return nil, fmt.Errorf("db not implemented")
	}

}
