package repository

import (
	"concert_pre-poster/internal/domain"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

type FirstVotingStagePsql struct {
	conn *sql.DB
}

func NewFirstVotingStagePsql(db *sql.DB) *FirstVotingStagePsql {
	return &FirstVotingStagePsql{
		conn: db,
	}
}

func (f *FirstVotingStagePsql) DoVote(idBillboard, idUser, idDate, maxTicketPrice int) (int, error) {

	tmp := f.conn.QueryRow(
		"INSERT INTO first_voting (id_billboard, id_user, id_date, max_ticket_price) VALUES ($1, $2, $3, $4) returning id",
		idBillboard, idUser, 1, maxTicketPrice,
	)
	var id int
	err := tmp.Scan(&id)

	if err != nil {
		log.Error(err)
		//return wrapErrorFromDB(err)
		return 0, err
	}
	return id, nil
}

func (f *FirstVotingStagePsql) DoVoteInBatch(idDates []int, idBillboard, idUser, maxTicketPrice int) error {
	for _, idDate := range idDates {
		_, err := f.conn.Exec(
			"INSERT INTO first_voting (id_billboard, id_user, id_date, max_ticket_price) VALUES ($1, $2, $3, $4)",
			idBillboard, idUser, idDate, maxTicketPrice,
		)
		if err != nil {
			log.Error(err)
			return wrapErrorFromDB(err)
		}
	}
	return nil
}

func (f *FirstVotingStagePsql) GetFirstVotingInfoForUser(userId int) (*[]domain.FirstVoting, error) {
	rows, err := f.conn.Query("SELECT * FROM first_voting WHERE id_user = ?", userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()
	var firstVotings []domain.FirstVoting
	for rows.Next() {
		firstVoting := domain.FirstVoting{}
		err = rows.Scan(&firstVoting.Id, &firstVoting.IdBillboard, &firstVoting.IdUser, &firstVoting.IdDate, &firstVoting.MaxTicketPrice)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		firstVotings = append(firstVotings, firstVoting)
	}

	return &firstVotings, nil
}

func (f *FirstVotingStagePsql) GetFirstVotingInfoForBillboard(billboardId int) (*[]domain.FirstVoting, error) {
	rows, err := f.conn.Query("SELECT * FROM first_voting WHERE id_billboard = ?", billboardId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()
	var firstVotings []domain.FirstVoting
	for rows.Next() {
		firstVoting := domain.FirstVoting{}
		err = rows.Scan(&firstVoting.Id, &firstVoting.IdBillboard, &firstVoting.IdUser, &firstVoting.IdDate, &firstVoting.MaxTicketPrice)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		firstVotings = append(firstVotings, firstVoting)
	}

	return &firstVotings, nil
}

func (f *FirstVotingStagePsql) AddDate(idBillboard int, date time.Time) (int, error) {
	tmp := f.conn.QueryRow(
		"INSERT INTO date (id_billboard, date )  VALUES ($1, $2) returning id",
		idBillboard, date,
	)
	var id int
	err := tmp.Scan(&id)

	if err != nil {
		log.Error(err)
		return 0, wrapErrorFromDB(err)
	}
	return id, nil
}

func (f *FirstVotingStagePsql) GetDateById(id int) (time.Time, error) {
	row := f.conn.QueryRow("select date from date where id = $1", id)

	var date time.Time

	err := row.Scan(&date)
	if err != nil {
		log.Error(err)
		return time.Time{}, err
	}
	return date, nil
}

func (f *FirstVotingStagePsql) AddDatesInBatch(idBillboard int, dates []time.Time) error {
	for _, date := range dates {
		_, err := f.conn.Exec(
			"INSERT INTO date (id_billboard, date)  VALUES ($1, $2)",
			idBillboard, date,
		)
		if err != nil {
			log.Error(err)
			return wrapErrorFromDB(err)
		}
	}
	return nil
}

func (f *FirstVotingStagePsql) GetMetrics(idBillboard int) (int, float64, error) {
	row := f.conn.QueryRow("select count(*), avg(max_ticket_price) from first_voting where id_billboard = $1", idBillboard)

	var count int
	var average_ticket_price float64

	err := row.Scan(&count, &average_ticket_price)
	if err != nil {
		log.Error("Error while scanning values", err)
		return 0, 0, err
	}
	return count, average_ticket_price, nil
}
