package api

import (
	"net/http"

	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "simple bank api",
		})
	})

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts", server.ListAccount)
	router.PUT("/accounts", server.UdateAccount)
	router.GET("/accounts/:id", server.GetAccount)
	router.DELETE("/accounts/:id", server.DeleteAccount)

	router.POST("/transfers", server.CreateTransfer)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
