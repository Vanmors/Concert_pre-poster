package repository

import (
	"concert_pre-poster/internal/domain"
	"database/sql"
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

func (f *FirstVotingStagePsql) DoVote(idDate, idBillboard, idUser, maxTicketPrice int) error {
	_, err := f.conn.Exec(
		"INSERT INTO first_voting (id_billboard, id_user, id_date, max_ticket_price) VALUES (?, ?, ?, ?)",
		idBillboard, idUser, idDate, maxTicketPrice,
	)
	if err != nil {
		return wrapErrorFromDB(err)
	}
	return nil
}

func (f *FirstVotingStagePsql) DoVoteInBatch(idDates []int, idBillboard, idUser, maxTicketPrice int) error {
	for _, idDate := range idDates {
		_, err := f.conn.Exec(
			"INSERT INTO first_voting (id_billboard, id_user, id_date, max_ticket_price) VALUES (?, ?, ?, ?)",
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