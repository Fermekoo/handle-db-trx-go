package api

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type RegisterRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	Fullname          string    `json:"fullname"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func userResponse(user db.User) UserResponse {
	return UserResponse{
		Fullname:          user.Fullname,
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) Register(ctx *gin.Context) {
	var request RegisterRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashed_password, err := utils.HashPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Fullname: request.Fullname,
		Username: request.Username,
		Email:    request.Email,
		Password: hashed_password,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := userResponse(user)

	ctx.JSON(http.StatusCreated, response)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func (server *Server) Login(ctx *gin.Context) {
	var request LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, request.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(request.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("email or password incorrect")))
		return
	}

	token, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("email or password incorrect")))
		return
	}

	userResponse := userResponse(user)

	response := LoginResponse{
		AccessToken: token,
		User:        userResponse,
	}

	ctx.JSON(http.StatusOK, response)

}
