package template

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisConns map[string]*redis.Client
type sqlConns map[string]*sql.DB

// Connection is a map to all data store connections used by a service.
type Connection struct {
	config ConfigMap
	redis  redisConns
	sql    sqlConns
}

// NewConnection creates a Connection from a ConfigMap.
func NewConnection(cm ConfigMap) *Connection {
	return &Connection{
		config: cm,
		redis:  make(redisConns),
		sql:    make(sqlConns),
	}
}

// GetRedis returns a redis connection specified by the namespace.
func (sc *Connection) GetRedis(namespace string) *redis.Client {
	return sc.redis[namespace]
}

// GetSQL returns a SQL connection specified by the namespace.
func (sc *Connection) GetSQL(namespace string) *sql.DB {
	return sc.sql[namespace]
}

// Connect to all databases that are configured.
func (sc *Connection) Connect() error {
	for ns, config := range sc.config {
		if config.ConfigType() == ConfigType.SQL {
			db, err := config.(SQLConnector).Connect()
			if err != nil {
				return err
			}
			sc.sql[ns] = db
		} else if config.ConfigType() == ConfigType.Redis {
			redis, err := config.(RedisConnector).Connect()
			if err != nil {
				return err
			}
			sc.redis[ns] = redis
		}
	}
	return nil
}

// Close all connections.
func (sc *Connection) Close() error {
	var err error
	for namespace := range sc.sql {
		db := sc.sql[namespace]
		if db == nil {
			continue
		}
		dbErr := db.Close()
		if dbErr == nil {
			delete(sc.sql, namespace)
		} else {
			err = errors.Wrap(dbErr, fmt.Sprintf("failed to close SQL connection: %s", namespace))
		}
	}
	for namespace := range sc.redis {
		client := sc.redis[namespace]
		if client == nil {
			continue
		}
		clientErr := client.Close()
		if clientErr == nil {
			delete(sc.redis, namespace)
		} else {
			err = errors.Wrap(clientErr, fmt.Sprintf("failed to close redis connection: %s", namespace))
		}
	}
	return err
}
