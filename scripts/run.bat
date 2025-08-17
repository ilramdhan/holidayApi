@echo off

REM Create data directory
if not exist "data" mkdir data

REM Set environment variables
set SERVER_HOST=localhost
set SERVER_PORT=8080
set DATABASE_PATH=./data/holidays.db
set MIGRATIONS_PATH=./migrations
set RATE_LIMIT_RPM=60
set RATE_LIMIT_BURST=10
set JWT_SECRET_KEY=super-secret-jwt-key-for-development-change-in-production
set JWT_ACCESS_TOKEN_TTL=15m
set JWT_REFRESH_TOKEN_TTL=168h
set ADMIN_API_KEY=admin-secret-key

REM Run the application
echo ==============================================
echo 🚀 Starting Holiday API Indonesia v2.0
echo ==============================================
echo Server: http://localhost:8080
echo Swagger: http://localhost:8080/swagger/index.html
echo Health: http://localhost:8080/health
echo.
echo 🔐 Admin Setup Required:
echo Please create admin user via database or API
echo See docs/DEPLOYMENT.md for secure setup
echo Default credentials have been removed for security
echo.
echo 📚 Quick Start:
echo 1. Login: POST /api/v1/auth/login
echo 2. Get holidays: GET /api/v1/holidays
echo 3. Admin operations: Use JWT token from login
echo ==============================================
echo.

go run cmd/server/main.go
