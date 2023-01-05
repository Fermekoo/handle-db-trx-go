migrateup:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/simplebank?sslmode=disable" -verbose down 

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

cleantestcache:
	go clean -testcache

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Fermekoo/handle-db-tx-go/db/sqlc Store

createmigrate:
	migrate create -ext sql -dir db/migrations -seq $(name)

postgres:
	docker run --name postgresdb -p 5432:5432 --network db-trx-go-network -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=root -d postgres

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb \
	--grpc-gateway_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9999 -r repl


.PHONY: migrateup migratedown sqlc test cleantestcache server mock createmigrate postgres proto evans