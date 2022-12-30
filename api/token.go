package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var request RenewTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refresh_payload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refresh_payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("refresh token is not valid")))
		return
	}

	if session.UserID != refresh_payload.UserID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("refresh token is not valid")))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("expired session")))
		return
	}

	access_token, payload, err := server.tokenMaker.CreateToken(session.UserID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := RenewTokenResponse{
		AccessToken:          access_token,
		AccessTokenExpiresAt: payload.ExpiredAt,
	}

	ctx.JSON(http.StatusOK, response)
}
