package repository

import (
	"concert_pre-poster/internal/domain"
	_ "concert_pre-poster/pkg/store/sqlstore"
	"database/sql"
	_ "time"

	log "github.com/sirupsen/logrus"
)

type UserPsql struct {
	conn *sql.DB
}

func NewUserPsql(db *sql.DB) *BillboardPsql {
	return &BillboardPsql{
		conn: db,
	}
}

func (u *UserPsql) GetUserByEmail(email string) (domain.User, error) {
	rows, err := u.conn.Query("SELECT * FROM user WHERE email = ?", email)

	if err != nil {
		log.Error(err)
		return domain.User{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return domain.User{}, sql.ErrNoRows
	}

	user := domain.User{}
	err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Hashed_password, &user.Phone_number, &user.Birthday)
	if err != nil {
		log.Error(err)
		return domain.User{}, err
	}

	return user, nil
}
