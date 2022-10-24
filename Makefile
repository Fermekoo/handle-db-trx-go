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

.PHONY: migrateup migratedown sqlc test cleantestcache server