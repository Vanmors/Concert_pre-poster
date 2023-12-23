package repository

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func GetBillboard_test(t *testing.T) {
	err := GetBillboard()
	assert.NoError(t, err)

}

