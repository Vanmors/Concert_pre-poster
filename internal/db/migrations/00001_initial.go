package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up00001Initial, Down00001Initial)
}

func Up00001Initial(_ context.Context, tx *sql.Tx) error {
	if err := createEnum(tx); err != nil {
		return err
	}
	if err := createUsersTable(tx); err != nil {
		return err
	}
	if err := createPerformerTable(tx); err != nil {
		return err
	}
	if err := createArticleTable(tx); err != nil {
		return err
	}
	if err := createBillboardTable(tx); err != nil {
		return err
	}
	if err := createDateTable(tx); err != nil {
		return err
	}
	if err := createFirstVotingTable(tx); err != nil {
		return err
	}
	if err := createSecondVotingTable(tx); err != nil {
		return err
	}
	if err := createRolesTable(tx); err != nil {
		return err
	}
	if err := createUserRolesTable(tx); err != nil {
		return err
	}
	if err := createResultSecondVotingTable(tx); err != nil {
		return err
	}
	if err := createCommentsTable(tx); err != nil {
		return err
	}
	if err := createPerformerBillboardTable(tx); err != nil {
		return err
	}
	return nil
}

func Down00001Initial(_ context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DROP TABLE IF EXISTS perfomer_billboard`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS comments`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS result_second_voting`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS user_roles`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS roles`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS second_voting`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS first_voting`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS date`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS billboard`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS article`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS performer`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TABLE IF EXISTS users`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DROP TYPE IF EXISTS comment_status`); err != nil {
		return err
	}

	return nil
}

func createEnum(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TYPE comment_status AS ENUM('ACCEPTED', 'HIDDEN', 'REJECTED');`)
	return err
}

func createUsersTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE users(
	  Id SERIAL PRIMARY KEY,
	  Name VARCHAR(50) NOT NULL,
	  Email VARCHAR(50) NOT NULL,
	  Hashed_password VARCHAR(255) NOT NULL,
	  Phone_number VARCHAR(20) NOT NULL,
	  Birthday TIMESTAMP NOT NULL
	);
`)
	return err
}

func createPerformerTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE performer(
	  Id SERIAL PRIMARY KEY,
	  Id_user INTEGER NOT NULL REFERENCES users(Id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Nickname VARCHAR(50) NOT NULL,
	  Genre VARCHAR(50) NOT NULL,
	  Description TEXT,
	  Verification BOOLEAN DEFAULT FALSE
	);
`)
	return err
}

func createArticleTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE article(
	  Id_article SERIAL PRIMARY KEY,
	  Id_performer INTEGER NOT NULL REFERENCES performer(id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Article TEXT NOT NULL
	);
`)
	return err
}

func createBillboardTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE billboard (
	  Id SERIAL PRIMARY KEY,
	  Relevance BOOLEAN DEFAULT TRUE,
	  Description TEXT,
	  City VARCHAR(100) NOT NULL,
	  Age_limit INTEGER
	);
`)
	return err
}

func createDateTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE date (
	  Id SERIAL PRIMARY KEY,
	  Id_billboard INTEGER NOT NULL REFERENCES billboard(id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Date TIMESTAMP
	);
`)
	return err
}

func createFirstVotingTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE first_voting (
	  Id SERIAL PRIMARY KEY,
	  Id_billboard INTEGER NOT NULL REFERENCES billboard(Id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Id_user INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Id_date INTEGER NOT NULL REFERENCES date(Id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Max_ticket_price INTEGER NOT NULL
	);
`)
	return err
}

func createSecondVotingTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE second_voting (
	  Id SERIAL PRIMARY KEY,
	  Id_billboard INTEGER NOT NULL REFERENCES billboard(Id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Id_user INTEGER NOT NULL,
	  Fan_vote BOOLEAN DEFAULT FALSE
	);
`)
	return err
}

func createRolesTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE  TABLE roles (
	  Id SERIAL PRIMARY KEY,
	  Name VARCHAR(30) NOT NULL
	);
`)
	return err
}

func createUserRolesTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE user_roles (
	  Id_user INTEGER NOT NULL REFERENCES users(Id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Id_role INTEGER NOT NULL REFERENCES roles(Id) ON UPDATE CASCADE ON DELETE CASCADE 
	);
`)
	return err
}

func createResultSecondVotingTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE result_second_voting(
		Id_billboard INTEGER NOT NULL REFERENCES billboard(Id) ON UPDATE CASCADE ON DELETE CASCADE,
		Final_date TIMESTAMP,
		Final_price INTEGER,
		Count_of_people INTEGER NOT NULL,
		Will_be_concert BOOLEAN NOT NULL
	);
`)
	return err
}

func createCommentsTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE comments(
	  Id SERIAL PRIMARY KEY,
	  Id_billboard INTEGER NOT NULL REFERENCES billboard(id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Id_user INTEGER NOT NULL REFERENCES users(Id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Comment TEXT NOT NULL,
	  Date TIMESTAMP NOT NULL,
	  Status comment_status
	);
`)
	return err
}

func createPerformerBillboardTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
	CREATE TABLE perfomer_billboard(
	  Id_performer INTEGER NOT NULL REFERENCES performer(id) ON UPDATE CASCADE ON DELETE CASCADE,
	  Id_billboard INTEGER NOT NULL REFERENCES billboard(id) ON UPDATE CASCADE ON DELETE CASCADE
	);
`)
	return err
}
