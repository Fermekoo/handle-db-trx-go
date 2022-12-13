migrateup:
	migrate -path db/migrations -database "postgresql://postgres:SNNCMEDT9wcqLsqr1kE3@simplebank.cd59x4i5kne2.ap-southeast-1.rds.amazonaws.com:5432/simplebank" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:SNNCMEDT9wcqLsqr1kE3@simplebank.cd59x4i5kne2.ap-southeast-1.rds.amazonaws.com:5432/simplebank" -verbose down 

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

.PHONY: migrateup migratedown sqlc test cleantestcache server mock createmigrate postgres