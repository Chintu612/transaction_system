# Makefile

# Load environment variables from .env file
include development.env
export $(shell sed 's/=.*//' development.env)

.PHONY: migrate-up
migrate-up:
	migrate -path db/migrations/ -database "$(DATABASE_URL)" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path db/migrations/ -database "$(DATABASE_URL)" -verbose down

.PHONY: migrate-create ## Create a DB migration files e.g `make migrate-create name=migration-name`
migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

.PHONY: all
all: migrate-up

.PHONY: clean
clean: migrate-down

