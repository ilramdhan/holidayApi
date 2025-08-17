#!/bin/bash

# Create data directory
mkdir -p data

# Set environment variables
export SERVER_HOST=localhost
export SERVER_PORT=8080
export DATABASE_PATH=./data/holidays.db
export MIGRATIONS_PATH=./migrations
export RATE_LIMIT_RPM=60
export RATE_LIMIT_BURST=10
export ADMIN_API_KEY=admin-secret-key

# Run the application
echo "Starting Holiday API Indonesia..."
echo "Server will be available at: http://localhost:8080"
echo "Swagger documentation: http://localhost:8080/swagger/index.html"
echo "Admin API Key: admin-secret-key"
echo ""

go run cmd/server/main.go
