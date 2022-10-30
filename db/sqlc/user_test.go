package db

import (
	"context"
	"testing"

	"github.com/Fermekoo/handle-db-tx-go/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Fullname: utils.RandomOwner(),
		Email:    utils.RandomEmail(),
		Username: utils.RandomOwner(),
		Password: utils.RandomString(64),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	get_user, err := testQueries.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, get_user)
	require.Equal(t, user.Fullname, get_user.Fullname)
	require.Equal(t, user.Email, get_user.Email)
	require.Equal(t, user.Username, get_user.Username)
	require.Equal(t, user.Password, get_user.Password)
}
