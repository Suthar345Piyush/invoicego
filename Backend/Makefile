.PHONY: help run build migrate-up migrate-down migrate-create sqlc-generate

help:
	@echo "Available commands:"
	@echo "  make run              - Run the application with air (hot reload)"
	@echo "  make build            - Build the application"
	@echo "  make migrate-up       - Run database migrations up"
	@echo "  make migrate-down     - Run database migrations down"
	@echo "  make migrate-create   - Create a new migration (usage: make migrate-create name=create_users)"
	@echo "  make sqlc-generate    - Generate Go code from SQL queries"
	@echo "  make docker-up        - Start PostgreSQL in Docker"
	@echo "  make docker-down      - Stop PostgreSQL in Docker"

run:
	air

build:
	go build -o bin/api cmd/api/main.go

migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/invoicego?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/invoicego?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

sqlc-generate:
	sqlc generate

docker-up:
	docker run --name invoice-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=invoice_generator -p 5432:5432 -d postgres:15-alpine

docker-down:
	docker stop invoice-postgres && docker rm invoice-postgres

test:
	go test -v ./...

lint:
	golangci-lint run



