package template

import (
	"context"
	"fmt"
	"net"

	redis "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	// Postgres driver
	_ "github.com/lib/pq"
)

// MySQLConfig contains parameters necessary to connect to a MySQL database.
type MySQLConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     int    `default:"3306"`
	User     string `required:"true"`
	Password string `required:"true"`
	DbName   string `envconfig:"DB_NAME" required:"true"`
}

// ConfigType of MySQLConfig is SQL.
func (c *MySQLConfig) ConfigType() ConfigTypeValue {
	return ConfigType.SQL
}

// Connect makes a SQL connection to the MySQL database.
func (c *MySQLConfig) Connect() (*sqlx.DB, error) {
	connPath := fmt.Sprintf("%v:%v@tcp4(%v:%v)/%v?parseTime=true", c.User, c.Password, c.Host, c.Port, c.DbName)
	db, err := sqlx.Open("mysql", connPath)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to MySQL")
	}
	return db, nil
}

// NetConfig contains the host and port parameters for a network listener.
type NetConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     string `default:"80"`
	Protocol string `default:"tcp"`
}

// ConfigType of NetConfig is NET.
func (c *NetConfig) ConfigType() ConfigTypeValue {
	return ConfigType.Net
}

// Address returns a string address with host and port.
func (c *NetConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Listen returns a net.Listener as specified by this NetConfig.
func (c *NetConfig) Listen() (net.Listener, error) {
	return net.Listen(c.Protocol, c.Address())
}

// PostgresConfig contains parameters necessary to connect to a Postgres database.
type PostgresConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     int    `default:"3306"`
	User     string `required:"true"`
	Password string `required:"true"`
	DbName   string `envconfig:"DB_NAME" required:"true"`
}

// ConfigType of PostgresConfig is SQL.
func (c *PostgresConfig) ConfigType() ConfigTypeValue {
	return ConfigType.SQL
}

// Connect makes a SQL connection to the Postgres database.
func (c *PostgresConfig) Connect() (*sqlx.DB, error) {
	connPath := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName,
	)
	db, err := sqlx.Open("postgres", connPath)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Postgres")
	}
	return db, nil
}

// RedisConfig contains parameters to connect to a redis database.
type RedisConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     string `default:"6379"`
	Password string
	DbNumber int `envconfig:"DB_NUMBER"`
}

// ConfigType of RedisConfig is REDIS.
func (c *RedisConfig) ConfigType() ConfigTypeValue {
	return ConfigType.Redis
}

// Connect makes a connection to the redis database.
func (c *RedisConfig) Connect() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Password: c.Password,
		DB:       c.DbNumber,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to redis")
	}
	return client, nil
}
