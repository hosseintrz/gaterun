package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hosseintrz/gaterun/pkg/api/util"
	"github.com/hosseintrz/gaterun/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func insertConsumer(ctx context.Context, consumer *Consumer) (id int64, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = database.ExecTx(db, func(tx *gorm.DB) error {
		exists, err := userExists(tx, consumer)
		if err != nil {
			return err
		}

		if exists {
			return util.InvalidRequestError("consumer with this username already exists", err)
		}

		err = db.Model(&Consumer{}).Create(consumer).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				err = util.NewHTTPError(http.StatusBadRequest, "consumer already exists", err)
			}
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	return consumer.ID, nil
}

func updateConsumer(ctx context.Context, id int64, consumer *Consumer) (*Consumer, error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	consumer.ID = id
	res := db.Model(&Consumer{}).Clauses(clause.Returning{}).Where(&Consumer{ID: id}).Updates(consumer)
	if res.RowsAffected == 0 {
		err = util.NewHTTPError(http.StatusNotFound, fmt.Sprintf("consumer with id : %d not found", id), res.Error)
		return nil, err
	}
	if err = res.Error; err != nil {
		return nil, err
	}

	return consumer, nil
}

func getConsumer(ctx context.Context, id int64) (consumer *Consumer, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Model(&Consumer{}).First(&consumer, &Consumer{ID: id}).Error
	if err != nil {
		return
	}

	return
}

func deleteConsumer(ctx context.Context, id int64) (consumer *Consumer, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Model(&Consumer{}).Clauses(clause.Returning{}).Delete(&consumer, id).Error
	if err != nil {
		return
	}

	return
}

func userExists(db *gorm.DB, consumer *Consumer) (exists bool, err error) {
	err = db.Raw(`
		SELECT EXISTS(SELECT 1 FROM consumers WHERE username = ? OR id = ?)
	`, consumer.Username, consumer.ID).Scan(&exists).Error

	if err != nil {
		return
	}

	return exists, nil
}
