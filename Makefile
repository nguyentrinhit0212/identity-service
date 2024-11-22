include .env
export $(shell sed 's/=.*//' .env)

MIGRATE_CMD=migrate -path ./db/migrations \
	-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

run:
	go run cmd/main.go

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

create-migration:
	$(MIGRATE_CMD) create -ext sql -dir ./db/migrations -seq $(NAME)