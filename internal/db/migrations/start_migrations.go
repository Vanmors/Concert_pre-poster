package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

func Start(reload bool, sq *sql.DB) {
	log.Info("Start migrations")
	if reload {
		log.Info("Rollback the database to the first version")
		err := goose.DownTo(sq, ".", 0)
		if err != nil {
			log.Error("Error during rollback", err.Error())
		}
	}
	log.Info("Update the version of migration")
	err := goose.Up(sq, ".")
	if err != nil {
		log.Error("error when trying to update migration version", err.Error())
	}
}
