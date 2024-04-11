.PHONY: migrate

include .env
export

MIGRATIONS_PATH := ./migrations

.PHONY: up down migrate logs

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

migrate:
	docker-compose run --rm migrator go run ./cmd/migrator --storage-path=$(PG_URL) --migrations-path=$(MIGRATIONS_PATH)
