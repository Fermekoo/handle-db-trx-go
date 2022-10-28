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

.PHONY: migrateup migratedown sqlc test cleantestcache server mock