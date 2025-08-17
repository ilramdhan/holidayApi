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
export JWT_SECRET_KEY=super-secret-jwt-key-for-development-change-in-production
export JWT_ACCESS_TOKEN_TTL=15m
export JWT_REFRESH_TOKEN_TTL=168h
export ADMIN_API_KEY=admin-secret-key

# Run the application
echo "=============================================="
echo "🚀 Starting Holiday API Indonesia v2.0"
echo "=============================================="
echo "Server: http://localhost:8080"
echo "Swagger: http://localhost:8080/swagger/index.html"
echo "Health: http://localhost:8080/health"
echo ""
echo "🔐 Default Admin Credentials:"
echo "Username: admin"
echo "Password: Admin123!"
echo "Role: super_admin"
echo ""
echo "📚 Quick Start:"
echo "1. Login: POST /api/v1/auth/login"
echo "2. Get holidays: GET /api/v1/holidays"
echo "3. Admin operations: Use JWT token from login"
echo "=============================================="
echo ""

go run cmd/server/main.go
