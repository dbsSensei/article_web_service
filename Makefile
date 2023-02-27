DB_URL=postgresql://postgres:password@localhost:5432/article_web_service?sslmode=disable

network:
	docker network create articlewebservicenetwork

postgres:
	docker run --name postgres --network articlewebservicenetwork -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres:14-alpine

postgresrestart:
	docker restart postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root article_web_service

dropdb:
	docker exec -it postgres dropdb article_web_service

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

serverdev:
	nodemon --exec go run main.go --signal SIGTERM

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/dbssensei/articlewebservice/db/sqlc Store

.PHONY: network postgres createdb dropdb migrateup migratedown sqlc test server mock
