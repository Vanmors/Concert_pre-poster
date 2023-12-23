package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBillboard(t *testing.T) {
	billboards, err := GetBillboard()
	for _, val := range billboards{
		t.Logf("%+v", val)
	}
	assert.NoError(t, err)
}
