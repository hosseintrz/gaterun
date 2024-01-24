package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/hosseintrz/gaterun/config/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ErrNilDatabase = errors.New("db is nil")
)

var once sync.Once
var db *gorm.DB

func InitDatabase(cfg models.DatabaseConfig) error {
	var err error

	once.Do(func() {
		dsn := getDSN(cfg)

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,        // Don't include params in the SQL log
				Colorful:                  false,       // Disable color
			},
		)

		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{
			Logger: newLogger,
		})

	})

	if err != nil {
		return err
		//		slog.Error("error initializing database -> %v\n", err)
	}

	return nil
}

func getDSN(cfg models.DatabaseConfig) string {
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
		tx.AddError(err)
		return err
	}

	return nil
}
