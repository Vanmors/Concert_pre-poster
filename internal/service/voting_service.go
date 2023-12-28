package service

import (
	"concert_pre-poster/internal/repository"
	"strconv"
	"time"
)

type Service struct {
	repos *repository.Repositories
}

func NewService(Repos *repository.Repositories) *Service {
	return &Service{
		repos: Repos,
	}
}

func (s *Service) Create_voting_service(idBillboard string, stringDates []string) error {
	var layout = "2006-01-02T15:04"

	for _, date := range stringDates {
		parsedTime, err := time.Parse(layout, date)
		if err != nil {
			return err
		}
		idBillboardInt, err := strconv.Atoi(idBillboard)
		if err != nil {
			return err
		}
		_, err = s.repos.FirstVotingStage.AddDate(idBillboardInt, parsedTime)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *Service) CalculateMetricsFirstVoting(idBillboard int) (int, float64, error) {
	countVotes, averagePrice, err := s.repos.FirstVotingStage.GetMetrics(idBillboard)
	if err != nil {
		return 0, 0, err
	}

	return countVotes, averagePrice, nil
}
