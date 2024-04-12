
include .env
export

up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: up

down:
	docker-compose down --remove-orphans
.PHONY: down


test:
	go test -v ./tests/extended && go test -v ./tests



