package service

import (
	"concert_pre-poster/internal/repository"
	"strconv"
	"time"
)


func Create_voting_service(idBillboard string, stringDates []string) error {
	repos, err := repository.NewRepositories("concert_pre-poster", "postgres", "password")
//	repos, err := repository.NewRepositories("ToDelete", "postgres", "password")
	if err != nil {
		return err
	}

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
		_, err = repos.FirstVotingStage.AddDate(idBillboardInt, parsedTime)
		if err != nil {
			return err
		}

	}
	return nil
}


func CalculateMetricsFirstVoting(idBillboard int) (int, float64, error){
	repos, err := repository.NewRepositories("concert_pre-poster", "postgres", "nav461")
	if err != nil {
		return 0, 0, err

	}

	countVotes, averagePrice, err := repos.FirstVotingStage.GetMetrics(idBillboard)
	if err != nil {
		return 0, 0, err
	}

	return countVotes, averagePrice, nil
}