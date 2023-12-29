package repository

import (
	"concert_pre-poster/internal/domain"
	"database/sql"
	"fmt"
	"time"
)

/*
1) внести голос за первый этап
2) внести голоса за первый этап - целой пачкой дат
3) получить информацию по своему голосу - по id фанату
4) получить информацию по id предафиши
*/

type FirstVotingStagePsql struct {
	conn *sql.DB
}

func NewFirstVotingStagePsql(db *sql.DB) *FirstVotingStagePsql {
	return &FirstVotingStagePsql{
		conn: db,
	}
}

func (f *FirstVotingStagePsql) DoVote(idBillboard, idUser, idDate, maxTicketPrice int) (int, error) {
	/*
		_, err := f.conn.Exec(
			"INSERT INTO first_voting (id_billboard, id_user, id_date, max_ticket_price) VALUES (?, ?, ?, ?)",
			idBillboard, idUser, idDate, maxTicketPrice,
		)

	*/
	tmp := f.conn.QueryRow(
		"INSERT INTO first_voting (id_billboard, id_user, id_date, max_ticket_price) VALUES ($1, $2, $3, $4) returning id",
		idBillboard, idUser, idDate, maxTicketPrice,
	)
	var id int
	err := tmp.Scan(&id)

	if err != nil {
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
			return wrapErrorFromDB(err)
		}
	}
	return nil
}

func (f *FirstVotingStagePsql) GetFirstVotingInfoForUser(userId int) (*[]domain.FirstVoting, error) {
	rows, err := f.conn.Query("SELECT * FROM first_voting WHERE id_user = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var firstVotings []domain.FirstVoting
	for rows.Next() {
		firstVoting := domain.FirstVoting{}
		err = rows.Scan(&firstVoting.Id, &firstVoting.IdBillboard, &firstVoting.IdUser, &firstVoting.IdDate, &firstVoting.MaxTicketPrice)
		if err != nil {
			return nil, err
		}
		firstVotings = append(firstVotings, firstVoting)
	}

	return &firstVotings, nil
}

func (f *FirstVotingStagePsql) GetFirstVotingInfoForBillboard(billboardId int) (*[]domain.FirstVoting, error) {
	rows, err := f.conn.Query("SELECT * FROM first_voting WHERE id_billboard = ?", billboardId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var firstVotings []domain.FirstVoting
	for rows.Next() {
		firstVoting := domain.FirstVoting{}
		err = rows.Scan(&firstVoting.Id, &firstVoting.IdBillboard, &firstVoting.IdUser, &firstVoting.IdDate, &firstVoting.MaxTicketPrice)
		if err != nil {
			return nil, err
		}
		firstVotings = append(firstVotings, firstVoting)
	}

	return &firstVotings, nil
}

func (f *FirstVotingStagePsql) AddDate(idBillboard int, date time.Time) (int, error) {
	/*
		_, err := f.conn.Exec(
			"INSERT INTO date (id_billboard, date ) VALUES (?, ?)",
			idBillboard, date,
		)

	*/
	tmp := f.conn.QueryRow(
		"INSERT INTO date (id_billboard, date )  VALUES ($1, $2) returning id",
		idBillboard, date,
	)
	var id int
	err := tmp.Scan(&id)

	if err != nil {
		return 0, wrapErrorFromDB(err)
	}
	return id, nil
}

func (f *FirstVotingStagePsql) GetDateById(id int) (time.Time, error) {
	row := f.conn.QueryRow("select date from date where id = $1", id)

	var date time.Time

	err := row.Scan(&date)
	if err != nil {
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
		fmt.Println("Ошибка при сканировании значений:", err)
		return 0, 0, err
	}
	return count, average_ticket_price, nil
}
