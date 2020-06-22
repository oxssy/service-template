package config

import (
	"context"
	"fmt"
	template "service-template"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// RedisConfig contains parameters to connect to a redis database.
type RedisConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     string `default:"6379"`
	Password string
	DbNumber int
}

// ConfigType of RedisConfig is REDIS.
func (c *RedisConfig) ConfigType() template.ConfigTypeValue {
	return template.ConfigType.Redis
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
