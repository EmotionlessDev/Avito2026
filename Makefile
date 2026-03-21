include .env
ENV_FILE = .env.local

up:
	docker compose --env-file $(ENV_FILE) up -d

build:
	docker compose build

logs:
	docker compose --env-file $(ENV_FILE) logs -f

down:
	docker compose --env-file $(ENV_FILE) down

migrate-up:
	migrate -path=./migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up

migrate-down:
	migrate -path=./migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" down

lint:
	golangci-lint run
