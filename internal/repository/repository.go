package repository

import (
	"concert_pre-poster/internal/domain"
	"concert_pre-poster/pkg/store/sqlstore"
)

type Billboard interface {
	GetBillboard() ([]domain.Billboard, error)
}

type PreliminaryResults interface {
}

type FinalResults interface {
}

type FirstVotingStage interface {
	DoVote(idDate, idBillboard, idUser, maxTicketPrice int) error
	DoVoteInBatch(idDates []int, idBillboard, idUser, maxTicketPrice int) error
	GetFirstVotingInfoForUser(userId int) (*[]domain.FirstVoting, error)
	GetFirstVotingInfoForBillboard(billboardId int) (*[]domain.FirstVoting, error)
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
