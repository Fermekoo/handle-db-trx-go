package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/Fermekoo/handle-db-tx-go/api"
	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/gapi"
	"github.com/Fermekoo/handle-db-tx-go/pb"
	"github.com/Fermekoo/handle-db-tx-go/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	runGRPCServer(config, store)

}

func runGinServer(config utils.Config, store db.Store) {
	server, _ := api.NewServer(config, store)

	err := server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterHandleDBServer(gRPCServer, server)
	reflection.Register(gRPCServer) // allow client to explore available procedure on the server

	listener, err := net.Listen("tcp", config.GRPCServerAddress)

	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())

	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

}
