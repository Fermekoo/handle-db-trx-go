package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/Fermekoo/handle-db-tx-go/api"
	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/gapi"
	"github.com/Fermekoo/handle-db-tx-go/pb"
	"github.com/Fermekoo/handle-db-tx-go/utils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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
	go runGatewayServer(config, store)
	runGRPCServer(config, store)

}

func runGinServer(config utils.Config, store db.Store) {
	server, _ := api.NewServer(config, store)

	err := server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}

func runGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterHandleDBHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)

	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start HTTP Gateway Server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http gateway server")
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
