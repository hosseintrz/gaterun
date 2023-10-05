package database

import (
	"context"
	"fmt"
	"sync"

	"errors"

	"github.com/hosseintrz/gaterun/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNilDatabase = errors.New("db is nil")
)

var once sync.Once
var db *gorm.DB

func InitDatabase(cfg config.DatabaseConfig) error {
	var err error

	once.Do(func() {
		dsn := getDSN(cfg)

		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	})

	if err != nil {
		return err
		//		slog.Error("error initializing database -> %v\n", err)
	}

	return nil
}

func getDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Tehran",
		cfg.Username,
		cfg.Password,
		cfg.DbName,
		cfg.Port,
		cfg.SslMode,
	)
}

func GetDB(ctx context.Context) (*gorm.DB, error) {
	if db == nil {
		return nil, ErrNilDatabase
	}

	return db.WithContext(ctx), nil
}

func ExecTx(db *gorm.DB, txFunc func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err := txFunc(tx)
	if err != nil {
		return err
	}

	return nil
}
