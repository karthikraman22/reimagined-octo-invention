package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHello(t *testing.T) {
	expected := Model{}
	actual, _ := GetModelByID("TEST")

	assert.Equal(t, actual, expected)
}
