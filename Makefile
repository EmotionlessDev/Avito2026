.PHONY: up down seed migrate logs

up:
	docker compose up --build

down:
	docker compose down

migrate:
	docker-compose run --rm migrate

seed:
	@echo "Seeding database..."
	@docker compose cp ./seed.sql db:/tmp/seed.sql
	@docker compose exec -T db psql -U postgres -d booking -f /tmp/seed.sql
	@echo "Seed completed"

logs:
	docker compose logs -f app

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  up       - Build and start the application"
	@echo "  down     - Stop and remove the application containers"
	@echo "  migrate  - Run database migrations"
	@echo "  seed     - Seed the database with initial data"
	@echo "  logs     - Follow the application logs"
