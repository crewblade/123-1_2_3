
include .env
export

up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: up

down:
	docker-compose down --remove-orphans
.PHONY: down


test:
	docker-compose -f docker-compose-test.yaml up --build --abort-on-container-exit && docker-compose -f docker-compose-test.yaml down --volumes



