package cassandra

import (
	"sync"

	"github.com/gocql/gocql"
	"github.com/hosseintrz/gaterun/config/models"
)

var cassandraOnce *sync.Once
var cassandraDB *gocql.Session

func InitDatabase(cfg models.DatabaseConfig) (err error) {
	cassandraOnce.Do(func() {
		cluster := gocql.NewCluster(cfg.Host)
		cluster.Keyspace = "gaterun"
		cluster.Port = cfg.Port
		cluster.Timeout = cfg.Timeout

		var session *gocql.Session
		session, err = cluster.CreateSession()
		if err == nil {
			cassandraDB = session
		}
	})

	if err != nil {
		return
	}

	return nil
}

func GetDB(cfg models.DatabaseConfig) (*gocql.Session, error) {
	if cassandraDB == nil {
		err := InitDatabase(cfg)
		if err != nil {
			return nil, err
		}
	}
	return cassandraDB, nil
}
