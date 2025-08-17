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
set ADMIN_API_KEY=admin-secret-key

REM Run the application
echo Starting Holiday API Indonesia...
echo Server will be available at: http://localhost:8080
echo Swagger documentation: http://localhost:8080/swagger/index.html
echo Admin API Key: admin-secret-key
echo.

go run cmd/server/main.go
