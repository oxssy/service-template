package config

import (
	"database/sql"
	"fmt"

	template "github.com/oxssy/service-template"

	_ "github.com/go-sql-driver/mysql" // mysql driver import
	"github.com/pkg/errors"
)

// MySQLConfig contains parameters necessary to connect to a MySQL database.
type MySQLConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     int    `default:"3306"`
	User     string `required:"true"`
	Password string `required:"true"`
	DbName   string `required:"true"`
}

// ConfigType of MySQLConfig is SQL.
func (c *MySQLConfig) ConfigType() template.ConfigTypeValue {
	return template.ConfigType.SQL
}

// Connect makes a SQL connection to the MySQL database.
func (c *MySQLConfig) Connect() (*sql.DB, error) {
	connPath := fmt.Sprintf("%v:%v@tcp4(%v:%v)/%v?parseTime=true", c.User, c.Password, c.Host, c.Port, c.DbName)
	db, err := sql.Open("mysql", connPath)
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to MySQL")
	}
	return db, nil
}
