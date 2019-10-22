package db

import (
	"database/sql"
	"fmt"

	"github.com/Wallruzz9114/bookey/config"
	_ "github.com/lib/pq" // postgres driver
)

// New ...
func New(config *config.Config) (*sql.DB, error) {
	dbInfo := fmt.Sprintf(
		"user=%s password=%s host=%v port=%v dbname=%s sslmode=disable",
		config.Db.Username, config.Db.Password, config.Db.Host, config.Db.Port, config.Db.DbName)

	return sql.Open("postgres", dbInfo)
}
