.PHONY: build run migrate test lint test.integration

build:
	docker-compose build avito-shop-service

run:
	docker-compose up avito-shop-service

env:
	source .env && echo "DB_USER: ${DB_USER}, DB_PASSWORD: ${DB_PASSWORD}, DB_NAME: ${DB_NAME}"

migrate:
	migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@localhost:5434/${DB_NAME}?sslmode=disable" up

test:
	go test -v ./...

lint:
	golangci-lint run

test.integration: 
	@echo "Starting PostgreSQL container..."
	docker run --rm --name test-postgres -e POSTGRES_USER=testuser -e POSTGRES_PASSWORD=testpassword -e POSTGRES_DB=testdb -p 5433:5432 -d postgres

	@echo "Waiting for PostgreSQL to start..."
	@until docker exec test-postgres pg_isready -U testuser -d testdb; do \
		sleep 2; \
	done

	@echo "Running database migrations..."
	migrate -path ./migrations -database "postgres://testuser:testpassword@localhost:5433/testdb?sslmode=disable" up

	@echo "Running Go tests..."
	go test -v -tags=integration ./tests

	docker stop test-postgres


