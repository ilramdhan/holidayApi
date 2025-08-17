.PHONY: build run test clean swagger deps migrate

# Build the application
build:
	go build -o bin/holidayapi cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml
	rm -f coverage.out coverage.html

# Create data directory
create-dirs:
	mkdir -p data

# Run with sample environment
dev: create-dirs
	SERVER_HOST=localhost \
	SERVER_PORT=8080 \
	DATABASE_PATH=./data/holidays.db \
	MIGRATIONS_PATH=./migrations \
	RATE_LIMIT_RPM=60 \
	RATE_LIMIT_BURST=10 \
	JWT_SECRET_KEY=super-secret-jwt-key-for-development \
	JWT_ACCESS_TOKEN_TTL=15m \
	JWT_REFRESH_TOKEN_TTL=168h \
	ADMIN_API_KEY=admin-secret-key \
	go run cmd/server/main.go

# Install swag tool
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Setup development environment
setup: deps install-swag create-dirs swagger

# Docker build
docker-build:
	docker build -t holidayapi .

# Docker run
docker-run:
	docker run -p 8080:8080 holidayapi
