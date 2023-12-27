package repository

import (
	"concert_pre-poster/pkg/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBillboard(t *testing.T) {
	db, err := sqlstore.NewClient("concert_pre-poster", "postgres", "password")
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
