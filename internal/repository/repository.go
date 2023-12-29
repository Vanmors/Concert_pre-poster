package repository

import (
	"concert_pre-poster/internal/domain"
	"concert_pre-poster/pkg/store/sqlstore"
	"time"
)

type Billboard interface {
	GetBillboard() ([]domain.Billboard, error)
	DeleteBilboardById(id int) error
	AddBillboard(relevance bool, description string, city string, age_limit int) (int, error)
	GetBillboardByID(id int) (domain.Billboard, error)
	GetBillboardAvailableDates(billboardId int) ([]*domain.Date, error)
}

type PreliminaryResults interface {
}

type FinalResults interface {
}

type FirstVotingStage interface {
	DoVote(idDate, idBillboard, idUser, maxTicketPrice int) (int, error)
	DoVoteInBatch(idDates []int, idBillboard, idUser, maxTicketPrice int) error
	GetFirstVotingInfoForUser(userId int) (*[]domain.FirstVoting, error)
	GetFirstVotingInfoForBillboard(billboardId int) (*[]domain.FirstVoting, error)
	AddDate(idBillboard int, date time.Time) (int, error)
	AddDatesInBatch(idBillboard int, dates []time.Time) error
	GetMetrics(idBillboard int) (int, float64, error)
	GetDateById(id int) (time.Time, error)
}

type Repositories struct {
	Billboard          Billboard
	PreliminaryResults PreliminaryResults
	FinalResults       FinalResults
	FirstVotingStage   FirstVotingStage
}

func NewRepositories(dbname, username, password string) (*Repositories, error) {
	db, err := sqlstore.NewClient(dbname, username, password)
	if err != nil {
		return nil, wrapErrorFromDB(err)
	}
	return &Repositories{
		Billboard:        NewBillboardPsql(db),
		FirstVotingStage: NewFirstVotingStagePsql(db),
	}, nil
}
