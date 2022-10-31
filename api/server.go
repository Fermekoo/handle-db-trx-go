package api

import (
	"fmt"
	"net/http"

	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/token"
	"github.com/Fermekoo/handle-db-tx-go/utils"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {

	token_maker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: token_maker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.SetupRouter()
	return server, nil
}

func (server *Server) SetupRouter() {
	router := gin.Default()
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

	router.POST("/register", server.Register)
	router.POST("/login", server.Login)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
