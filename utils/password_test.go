package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password, err := HashPassword(RandomString(100))
	require.NoError(t, err)
	require.NotEmpty(t, password)
}

func TestCheckPassword(t *testing.T) {
	plain_password := RandomString(100)
	hashPassword, err := HashPassword(plain_password)
	require.NoError(t, err)

	err = CheckPassword(plain_password, hashPassword)
	require.NoError(t, err)

	wrong_password := RandomString(6)
	err = CheckPassword(wrong_password, hashPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
