package repository

import (
	"concert_pre-poster/pkg/store/sqlstore"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/spf13/viper"
)

func prepareRepo(t *testing.T) (*FirstVotingStagePsql, *BillboardPsql) {
	viper.SetConfigFile("../../config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Получаем значения из конфигурации
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")

	db, err := sqlstore.NewClient(dbname, username, password)
	assert.NoError(t, err)
	return NewFirstVotingStagePsql(db), NewBillboardPsql(db)
}

func prepareDBData(t *testing.T, billRepo *BillboardPsql, fr *FirstVotingStagePsql) (idDate int, idBillBoard int) {
	idBillBoard, err := billRepo.AddBillboard(true, "ATL Concer", "Saint-P", 18)
	assert.NoError(t, err)
	idDate, err = fr.AddDate(idBillBoard, time.Unix(20000, 312312))
	assert.NoError(t, err)
	return idDate, idBillBoard
}

func TestDoVote(t *testing.T) {
	firstVoteStageRepo, billRepo := prepareRepo(t)
	idDate, idBillBoard := prepareDBData(t, billRepo, firstVoteStageRepo)
	_, _ = idDate, idBillBoard
	id, err := firstVoteStageRepo.DoVote(idBillBoard, 1, idDate, 200)
	assert.NoError(t, err)
	t.Log("id: ", id)
}
