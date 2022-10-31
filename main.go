package main

import (
	"database/sql"
	"log"

	"github.com/Fermekoo/handle-db-tx-go/api"
	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/utils"

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
	server, _ := api.NewServer(config, store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
