package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewClient(dbname, username, password string) (*sql.DB, error) {
	connStr := "user=postgres password=mypass dbname=productdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
