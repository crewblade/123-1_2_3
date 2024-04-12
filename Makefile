.PHONY: migrate

include .env
export

MIGRATIONS_PATH := ./migrations

up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: up

down:
	docker-compose down --remove-orphans
.PHONY: down


migrate-up:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	echo "y" | migrate -path migrations -database '$(PG_URL)?sslmode=disable' down
.PHONY: migrate-down

test:
	go test -v ./tests



