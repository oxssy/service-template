package template

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
)

// A ConfigTypeValue is just a string.
type ConfigTypeValue string

// ConfigType consists of a few commonly used config types.
var ConfigType = struct {
	Net   ConfigTypeValue
	Redis ConfigTypeValue
	SQL   ConfigTypeValue
}{
	"NET",
	"REDIS",
	"SQL",
}

// RedisConnector is an interface for RedisConfig.
type RedisConnector interface {
	Connect() (*redis.Client, error)
}

// SQLConnector is an interface for all SQL type configs, such as MySQLConfig.
type SQLConnector interface {
	Connect() (*sqlx.DB, error)
}

// A Config is typically a struct containing the config values.
// For it to be loadable from environment, please use struct tags as specified in:
// github.com/kelseyhightower/envconfig
type Config interface {
	ConfigType() ConfigTypeValue
}

// ConfigMap is a map from namespace to Config.
// Use the Get and Set methods to access Configs in a ConfigMap.
type ConfigMap map[string]Config

// NewConfigMap returns a ConfigMap.
func NewConfigMap() ConfigMap {
	return make(ConfigMap)
}

// Exists returns true if namespace exists in this service config.
func (cm ConfigMap) Exists(namespace string) bool {
	_, ok := cm[namespace]
	return ok
}

// Set a config Params into a namespace.
// If the namespace already exists, panic.
func (cm ConfigMap) Set(namespace string, config Config) {
	if cm.Exists(namespace) {
		panic(fmt.Sprintf("Overwriting config namespace: %s", namespace))
	}
	cm[namespace] = config
}

// Get a config Params given the namespace.
// If the namespace does not exist, panic.
func (cm ConfigMap) Get(namespace string) Config {
	if !cm.Exists(namespace) {
		panic(fmt.Sprintf("Missing config namespace: %s", namespace))
	}
	return cm[namespace]
}

// Load the config values from the environment.
func (cm ConfigMap) Load() error {
	for ns, config := range cm {
		err := envconfig.Process(ns, config)
		if err != nil {
			return err
		}
	}
	return nil
}
