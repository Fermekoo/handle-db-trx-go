package gapi

import (
	"fmt"

	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/pb"
	"github.com/Fermekoo/handle-db-tx-go/token"
	"github.com/Fermekoo/handle-db-tx-go/utils"
)

type Server struct {
	pb.UnimplementedHandleDBServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server
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

	return server, nil
}
