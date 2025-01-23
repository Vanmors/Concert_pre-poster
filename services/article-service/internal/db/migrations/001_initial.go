package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up001Initial, Down001Initial)
}

func Up001Initial(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`
	CREATE TABLE article (
		id_article SERIAL PRIMARY KEY,
		id_performer INTEGER NOT NULL,
		article TEXT NOT NULL
	);
`); err != nil {
		return err
	}
	return nil
}

func Down001Initial(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS article`); err != nil {
		return err
	}
	return nil
}
