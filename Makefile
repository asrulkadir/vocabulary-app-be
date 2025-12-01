# Load DB_URL from .env file
include .env
export

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

migrate-status:
	migrate -path ./migrations -database "$(DATABASE_URL)" version

help:
	@echo "Available commands:"
	@echo "  make migrate-up      - Run all pending migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-status  - Show current migration version"
