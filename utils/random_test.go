package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	var min int64 = 10
	var max int64 = 5000
	randomInt := RandomInt(min, max)

	require.GreaterOrEqual(t, randomInt, min)
	require.LessOrEqual(t, randomInt, max)
}

func TestRandomString(t *testing.T) {
	var string_length int = 20

	randomString := RandomString(string_length)

	require.NotEmpty(t, randomString)
	require.Len(t, randomString, 20)
}

func TestRandomOwner(t *testing.T) {
	randomOwner := RandomOwner()

	require.NotEmpty(t, randomOwner)
}

func TestRandomMoney(t *testing.T) {
	randomMoney := RandomMoney()

	require.NotEmpty(t, randomMoney)
}

func TestRandomCurrency(t *testing.T) {
	randomCurrency := RandomCurrency()

	require.NotEmpty(t, randomCurrency)
}
