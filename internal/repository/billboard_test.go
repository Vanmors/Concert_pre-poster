package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetBillboard(t *testing.T) {
	err := GetBillboard()
	assert.NoError(t, err)
}
