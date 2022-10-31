package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Fermekoo/handle-db-tx-go/db/mock"
	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqRegisterParamsmatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqRegisterParamsmatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	hashed_password := arg.Password
	err := utils.CheckPassword(e.password, hashed_password)
	if err != nil {
		return false
	}

	e.arg.Password = arg.Password
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqRegisterParamsmatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqRegisterParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqRegisterParamsmatcher{arg, password}
}

func createRandomUser(t *testing.T) (user db.User, password string) {

	password = utils.RandomString(10)
	hashed_password, _ := utils.HashPassword(password)
	user = db.User{
		ID:       utils.RandomInt(1, 1000),
		Fullname: utils.RandomOwner(),
		Username: utils.RandomOwner(),
		Email:    utils.RandomEmail(),
		Password: hashed_password,
	}

	return user, password
}

func TestRegisterAPI(t *testing.T) {
	user, password := createRandomUser(t)
	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"fullname": user.Fullname,
				"username": user.Username,
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					Fullname: user.Fullname,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqRegisterParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/register"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var getUser db.User

	err = json.Unmarshal(data, &getUser)
	require.NoError(t, err)
	require.Equal(t, user.Fullname, getUser.Fullname)
	require.Equal(t, user.Username, getUser.Username)
	require.Equal(t, user.Email, getUser.Email)
}
