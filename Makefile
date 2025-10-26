include .env
export

DIR=internal/db/migrations
CONN_STRING=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)

server:
	go run ./cmd/api

sqlc:
	sqlc generate

migrate-create:
	migrate create -ext sql -dir $(DIR) -seq $(NAME)

migrate-up:
	migrate -path $(DIR) -database $(CONN_STRING) up

migrate-down:
	migrate -path $(DIR) -database $(CONN_STRING) down 1

migrate-force:
	migrate -path $(DIR) -database $(CONN_STRING) force $(VERSION)

.PHONY: migrate-create migrate-up migrate-down migrate-force server sqlc