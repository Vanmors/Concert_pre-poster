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
