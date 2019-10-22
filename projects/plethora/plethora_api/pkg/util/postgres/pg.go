package postgres

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	// DB adapter
	_ "github.com/lib/pq"
)

// databaseLogger ...
type databaseLogger struct{}

// BeforeQuery ...
func (logger databaseLogger) BeforeQuery(pg *pg.QueryEvent) {}

// AfterQuery ...
func (logger databaseLogger) AfterQuery(query *pg.QueryEvent) {
	log.Printf(query.FormattedQuery())
}

// New creates new database connection to a postgres database
func New(psn string, timeout int, enableLog bool) (*pg.DB, error) {
	u, err := pg.ParseURL(psn)

	if err != nil {
		return nil, err
	}

	db := pg.Connect(u)

	_, err = db.Exec("SELECT 1")

	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if enableLog {
		db.AddQueryHook(databaseLogger{})
	}

	return db, nil
}
