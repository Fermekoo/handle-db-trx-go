package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/token"
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR IDR"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {
	var request CreateAccountRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	auth_payload := ctx.MustGet(authorization_payload_key).(*token.Payload)
	arg := db.CreateAccountParams{
		UserID:   auth_payload.UserID,
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  0,
	}
	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var request GetAccountRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	auth_payload := ctx.MustGet(authorization_payload_key).(*token.Payload)

	account, err := server.store.GetAccount(ctx, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if account.UserID != auth_payload.UserID {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("account not found")))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) ListAccount(ctx *gin.Context) {
	var request ListAccountRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	auth_payload := ctx.MustGet(authorization_payload_key).(*token.Payload)
	arg := db.ListAccountsParams{
		UserID: auth_payload.UserID,
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}
	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

type UpdateAccountRequest struct {
	ID       int64  `json:"id,string" binding:"required"`
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR IDR"`
}

func (server *Server) UdateAccount(ctx *gin.Context) {
	var request UpdateAccountRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, request.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:       request.ID,
		Owner:    request.Owner,
		Currency: request.Currency,
		Balance:  account.Balance,
	}

	account, err = server.store.UpdateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type DeleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) DeleteAccount(ctx *gin.Context) {
	var request DeleteAccountRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, request.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
