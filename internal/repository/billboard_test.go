package repository

import (
	"concert_pre-poster/pkg/store/sqlstore"

	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetBillboard(t *testing.T) {

	viper.SetConfigFile("../../config/config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Получаем значения из конфигурации
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")

	db, err := sqlstore.NewClient(dbname, username, password)

	t.Logf("%+v", err)
	billboardRepo := NewBillboardPsql(db)
	billboards, err := billboardRepo.GetBillboard()
	for _, val := range billboards {
		t.Logf("%+v", val)
	}
	assert.NoError(t, err)
}

func TestBillboardPsql_GetBillboardAvailableDates(t *testing.T) {
	_, billRepo := prepareRepo(t)
	dates, err := billRepo.GetBillboardAvailableDates(1)
	assert.NoError(t, err)
	for _, val := range dates {
		t.Logf("%+v", val)
	}
}
