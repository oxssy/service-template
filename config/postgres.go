package config

import (
	"database/sql"
	"fmt"

	template "github.com/oxssy/service-template"

	_ "github.com/lib/pq" // Postgres sql driver
	"github.com/pkg/errors"
)

// PostgresConfig contains parameters necessary to connect to a Postgres database.
type PostgresConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     int    `default:"3306"`
	User     string `required:"true"`
	Password string `required:"true"`
	DbName   string `required:"true"`
}

// ConfigType of PostgresConfig is SQL.
func (c *PostgresConfig) ConfigType() template.ConfigTypeValue {
	return template.ConfigType.SQL
}

// Connect makes a SQL connection to the Postgres database.
func (c *PostgresConfig) Connect() (*sql.DB, error) {
	connPath := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName,
	)
	db, err := sql.Open("postgres", connPath)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Postgres")
	}
	return db, nil
}
