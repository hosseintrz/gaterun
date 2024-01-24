package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	config "github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/api/admin/models"
	"github.com/hosseintrz/gaterun/pkg/api/util"
	"github.com/hosseintrz/gaterun/pkg/auth"
	"github.com/hosseintrz/gaterun/pkg/cache/redis"
	"github.com/hosseintrz/gaterun/pkg/database"

	"github.com/hosseintrz/gaterun/pkg/database/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InsertConsumer(ctx context.Context, consumer *models.Consumer) (id int64, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = postgres.ExecTx(db, func(tx *gorm.DB) error {
		exists, err := UserExists(tx, consumer)
		if err != nil {
			return err
		}

		if exists {
			return util.InvalidRequestError("consumer with this username already exists", err)
		}

		err = db.Model(&models.Consumer{}).Create(consumer).Error
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

func UpdateConsumer(ctx context.Context, id int64, consumer *models.Consumer) (*models.Consumer, error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return nil, err
	}

	consumer.ID = id
	res := db.Model(&models.Consumer{}).Clauses(clause.Returning{}).Where(&models.Consumer{ID: id}).Updates(consumer)
	if res.RowsAffected == 0 {
		err = util.NewHTTPError(http.StatusNotFound, fmt.Sprintf("consumer with id : %d not found", id), res.Error)
		return nil, err
	}
	if err = res.Error; err != nil {
		return nil, err
	}

	return consumer, nil
}

func GetConsumer(ctx context.Context, id int64) (consumer *models.Consumer, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Model(&models.Consumer{}).First(&consumer, &models.Consumer{ID: id}).Error
	if err != nil {
		return
	}

	return
}

func DeleteConsumer(ctx context.Context, id int64) (consumer *models.Consumer, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Model(&models.Consumer{}).Clauses(clause.Returning{}).Delete(&consumer, id).Error
	if err != nil {
		return
	}

	return
}

func UserExists(db *gorm.DB, consumer *models.Consumer) (exists bool, err error) {
	err = db.Raw(`
		SELECT EXISTS(SELECT 1 FROM consumers WHERE username = ? OR id = ?)
	`, consumer.Username, consumer.ID).Scan(&exists).Error

	if err != nil {
		return
	}

	return exists, nil
}

func GenerateApiKey(ctx context.Context, consumerId int64) (apiKey *auth.ApiKey, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	var consumer models.Consumer
	err = db.Table("consumers").Take(&consumer, consumerId).Error
	if err != nil {
		return
	}

	token := uuid.New().String()

	apiKey = &auth.ApiKey{
		Key:              token,
		CreatedAt:        time.Now(),
		ConsumerId:       consumer.ID,
		ConsumerCustomId: consumer.CustomID,
		ConsumerUsername: consumer.Username,
		TTL:              1000,
	}

	jsonVal, err := json.Marshal(apiKey)
	if err != nil {
		return
	}

	redisClient := redis.GetClient(ctx)
	err = redisClient.Set(ctx, token, string(jsonVal), time.Duration(apiKey.TTL*int64(time.Second))).Err()
	if err != nil {
		return
	}

	return
}

func InsertEndpoint(ctx context.Context, req *models.EndpointRequestDTO) (insertedId int64, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = postgres.ExecTx(db, func(tx *gorm.DB) (err error) {
		insertedId, err = InsertEndpointTx(ctx, tx, req)
		if err != nil {
			return
		}

		return nil
	})
	if err != nil {
		return
	}

	return
}

func InsertEndpointTx(ctx context.Context, tx *gorm.DB, req *models.EndpointRequestDTO) (insertedId int64, err error) {
	err = tx.Table("endpoint").Omit("backends").Create(req).Error
	//err = tx.Table("endpoint").Create(&req).Error
	if err != nil {
		return
	}

	insertedId = req.ID
	mapping := make([]models.EndpointBackend, len(req.Backends))
	for i, idStr := range req.Backends {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return 0, err
		}
		mapping[i] = models.EndpointBackend{
			EndpointId: insertedId,
			BackendId:  int64(id),
		}
	}

	err = tx.Table("endpoint_backends").CreateInBatches(mapping, 10).Error
	if err != nil {
		return
	}

	return insertedId, nil
}

func FetchAllEndpoints(ctx context.Context) (endpoints []config.EndpointConfig, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Table("endpoints").Scan(&endpoints).Error
	if err != nil {
		return
	}

	return endpoints, nil
}

func FetchEndpoint(ctx context.Context, id int64) (endpoint *config.EndpointConfig, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	return FetchEndpointTx(ctx, db, id)
}

func FetchEndpointTx(ctx context.Context, tx *gorm.DB, id int64) (endpoint *config.EndpointConfig, err error) {
	endpoint = &config.EndpointConfig{}

	err = tx.Table("endpoint").Where("id", id).Take(&endpoint).Error
	endpoint.Timeout *= time.Second
	endpoint.CacheTTL *= time.Second

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = util.NotFoundError("endpoint not found", err)
		}
		return
	}

	backends := []*config.BackendConfig{}
	err = tx.Table("backend").
		Joins("LEFT JOIN endpoint_backends as eb on backend.id = eb.backend_id").
		Where("eb.endpoint_id", endpoint.ID).
		Select("backend.*").
		Scan(&backends).
		Error

	// rows, err := tx.Raw(`
	// 		SELECT id, host, method, url_pattern, allow_list, deny_list, mapping, encoding, timeout, decoder_factory, headers_to_pass, concurrent_calls
	// 		FROM backend
	// 		Where id IN (SELECT backend_id FROM endpoint_backends WHERE endpoint_id = $1)
	// 	`, endpoint.ID).Rows()

	if err != nil {
		return
	}

	endpoint.Backends = backends
	// defer rows.Close()

	// for rows.Next() {
	// 	var backend config.BackendConfig
	// 	var allowList, denyList, mapping string

	// 	err = rows.Scan(
	// 		&backend.ID,
	// 		&backend.Host,
	// 		&backend.Method,
	// 		&backend.URLPattern,
	// 		&allowList,
	// 		&denyList,
	// 		&mapping,
	// 		&backend.Encoding,
	// 		&backend.Timeout,
	// 		&backend.DecoderFactory,
	// 		&backend.HeadersToPass,
	// 		&backend.ConcurrentCalls,
	// 	)
	// 	if err != nil {
	// 		return
	// 	}

	// 	err = json.Unmarshal([]byte(mapping), &backend.Mapping)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	endpoint.Backends = append(endpoint.Backends, &backend)
	// }

	return endpoint, nil
}

func DeleteEndpoint(ctx context.Context, id int64) (endpoint *config.EndpointConfig, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Table("endpoints").Clauses(clause.Returning{}).Delete(&endpoint, id).Error
	if err != nil {
		return
	}

	return endpoint, nil
}

func InsertBackend(ctx context.Context, backend *config.BackendConfig) (insertedId int64, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Table("backend").Create(backend).Error
	if err != nil {
		return
	}

	return backend.ID, nil
}

func FetchAllBackends(ctx context.Context) (backends []*config.BackendConfig, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Table("backend").Scan(&backends).Error
	if err != nil {
		return
	}

	return backends, nil
}

func FetchBackend(ctx context.Context, id int64) (backend *config.BackendConfig, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Table("backend").Take(&backend, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = util.NotFoundError("backend not found", err)
		}
		return
	}

	return backend, nil
}

func DeleteBackend(ctx context.Context, id int64) (backend *config.BackendConfig, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	err = db.Table("backend").Clauses(clause.Returning{}).Delete(&backend, id).Error
	if err != nil {
		return
	}

	return backend, nil
}

func UpdateBackend(ctx context.Context, id int64, backend *config.BackendConfig) (insertedId int64, err error) {
	db, err := database.GetDB(ctx)
	if err != nil {
		return
	}

	backend.ID = id
	res := db.Model(&config.BackendConfig{}).Clauses(clause.Returning{}).Where(&config.BackendConfig{ID: id}).Updates(backend)
	if res.RowsAffected == 0 {
		err = util.NewHTTPError(http.StatusNotFound, fmt.Sprintf("consumer with id : %d not found", id), res.Error)
		return
	}
	if err = res.Error; err != nil {
		return
	}

	return backend.ID, nil
}

// func assembleConfig(ctx context.Context, tx *gorm.DB) (cfg *config.ServiceConfig, err error) {
// 	cfg, err = persistence.LoadConfigFromDb(ctx)
// 	if err != nil {
// 		return
// 	}

// 	endpointIds := []int64{}
// 	err = tx.Table("endpoint").Select("id").Scan(&endpointIds).Error
// 	if err != nil {
// 		return
// 	}

// 	for _, id := range endpointIds {
// 		endpoint, err := fetchEndpointTx(ctx, tx, id)
// 		if err != nil {
// 			return nil, err
// 		}

// 		cfg.Endpoints = append(cfg.Endpoints, endpoint)
// 	}

// 	return
// }
