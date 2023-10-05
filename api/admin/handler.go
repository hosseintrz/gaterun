package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hosseintrz/gaterun/api/util"
	"github.com/hosseintrz/gaterun/internal/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func insertConsumer(ctx context.Context, consumer *Consumer) (id int64, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Model(&Consumer{}).Create(consumer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			err = util.NewHTTPError(http.StatusBadRequest, "consumer with this username exists", err)
		}
		return
	}

	return consumer.ID, nil
}

func updateConsumer(ctx context.Context, id int64, consumer *Consumer) (err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	res := db.Model(&Consumer{}).Where(&Consumer{ID: id}).Updates(consumer)
	if res.RowsAffected == 0 {
		err = util.NewHTTPError(http.StatusNotFound, fmt.Sprintf("consumer with id : %d not found", id), res.Error)
		return
	}
	if err = res.Error; err != nil {
		return
	}

	return nil
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
