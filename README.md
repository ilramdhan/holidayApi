# 🎉 Holiday API Indonesia 🇮🇩

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/Swagger-API%20Docs-orange?style=for-the-badge&logo=swagger&logoColor=white" alt="Swagger">
  <img src="https://img.shields.io/badge/REST-API-blue?style=for-the-badge" alt="REST API">
</p>

<p align="center">
  <b>Production-ready REST API for Indonesian National Holidays & Joint Leave Days</b><br>
  Based on official SKB 3 Menteri regulations | JWT Authentication | Rate Limiting | Audit Logging
</p>

<p align="center">
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-api-documentation">API Docs</a> •
  <a href="#-use-cases">Use Cases</a> •
  <a href="#-sdk--clients">SDK</a> •
  <a href="#-self-hosting">Self-Hosting</a>
</p>

---

## 🚀 Live Demo

**Public API Endpoint:** `https://api.holidayapi.id/v1` *(Coming Soon - Currently Self-Host)*

```bash
# Try it now with curl
curl -X GET "https://api.holidayapi.id/v1/holidays/year/2024" \
  -H "Accept: application/json"
```

---

## ✨ Features

### 🔐 **Authentication & Security**
- ✅ **JWT Authentication** with access & refresh tokens
- ✅ **Role-Based Access Control (RBAC)** - Super Admin & Admin roles
- ✅ **User Management** - Registration, login, profile management
- ✅ **Password Security** - Bcrypt hashing with password policy
- ✅ **Comprehensive Audit Logging** - Track all user actions
- ✅ **Security Headers** - XSS, CSRF, Content Security Policy protection
- ✅ **Input Sanitization** - Protection from injection attacks
- ✅ **Enhanced Rate Limiting** - Per-user rate limiting (60 RPM public, authenticated users higher)

### 🚀 **API Features**
- ✅ **REST API** with versioning (v1)
- ✅ **CRUD operations** for admin (JWT protected)
- ✅ **Filter by type** - National holidays, joint leave days, or both
- ✅ **Filter by period** - Year, month, or specific day queries
- ✅ **Swagger documentation** with interactive testing
- ✅ **SQLite database** (pure Go, no CGO required)
- ✅ **Comprehensive logging** with structured format
- ✅ **Input validation** with custom validators
- ✅ **Docker support** with production-ready configuration
- ✅ **Unit & Integration tests** with mocking

---

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [API Documentation](#-api-documentation)
- [Use Cases](#-use-cases)
- [Who is This For](#-who-is-this-for)
- [Self-Hosting Guide](#-self-hosting)
- [SDK & Clients](#-sdk--clients)
- [Integrations](#-integrations)
- [Configuration](#-configuration)
- [Docker Deployment](#-docker-deployment)

---

## 🚀 Quick Start

### Prerequisites
- Go 1.23 or higher
- Git
- Docker (optional)

### Installation & Running

#### Option 1: Using Docker Compose (Recommended)
```bash
# Clone the repository
git clone https://github.com/ilramdhan/holidayapi.git
cd holidayapi

# Start with Docker Compose
docker-compose up --build

# API is now running at http://localhost:8080
```

#### Option 2: Using Go
```bash
# Clone the repository
git clone https://github.com/ilramdhan/holidayapi.git
cd holidayapi

# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go
```

#### Option 3: Install as Go Package
```bash
go get github.com/ilramdhan/holidayapi
```

### Verify Installation
```bash
# Health check
curl http://localhost:8080/health

# Get API documentation
curl http://localhost:8080/swagger/index.html
```

---

## 📚 API Documentation

### Interactive Documentation
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI Spec**: http://localhost:8080/swagger/doc.json

### Base URL
```
http://localhost:8080/api/v1
```

### 🔓 Public Endpoints (No Authentication Required)

| Endpoint | Description | Example |
|----------|-------------|---------|
| `GET /api/v1/holidays` | Get all holidays with filters | [Try it](#) |
| `GET /api/v1/holidays/year/{year}` | Get holidays by year | `/holidays/year/2024` |
| `GET /api/v1/holidays/month/{year}/{month}` | Get holidays by month | `/holidays/month/2024/1` |
| `GET /api/v1/holidays/today` | Get today's holiday | `/holidays/today` |
| `GET /api/v1/holidays/this-year` | Get current year holidays | `/holidays/this-year` |
| `GET /api/v1/holidays/upcoming` | Get upcoming holidays | `/holidays/upcoming` |
| `GET /health` | Health check | `/health` |

### 🔐 Authentication Endpoints

| Endpoint | Description | Auth Required |
|----------|-------------|---------------|
| `POST /api/v1/auth/login` | User login (get JWT tokens) | No |
| `POST /api/v1/auth/refresh` | Refresh access token | Refresh Token |
| `GET /api/v1/auth/profile` | Get user profile | JWT |
| `POST /api/v1/auth/change-password` | Change password | JWT |

### 👑 Admin Endpoints (JWT Required)

| Endpoint | Description | Role |
|----------|-------------|------|
| `POST /api/v1/admin/holidays` | Create new holiday | Admin/Super Admin |
| `PUT /api/v1/admin/holidays/{id}` | Update holiday | Admin/Super Admin |
| `DELETE /api/v1/admin/holidays/{id}` | Delete holiday | Admin/Super Admin |
| `GET /api/v1/admin/audit-logs` | View all audit logs | Super Admin |

---

## 💡 Use Cases

### 1. **HR & Payroll Systems**
Automatically calculate working days, leave balances, and payroll based on official Indonesian holidays.

```python
# Example: Calculate working days in a month
import requests

response = requests.get("https://api.holidayapi.id/v1/holidays/month/2024/1")
holidays = response.json()["data"]

working_days = total_days - len(holidays)
```

### 2. **Calendar & Scheduling Apps**
Display Indonesian holidays in calendar applications and scheduling tools.

### 3. **E-commerce & Retail**
Plan promotions and sales around holiday periods. Identify peak shopping seasons.

### 4. **Logistics & Delivery**
Optimize delivery schedules by accounting for non-working days and joint leave periods.

### 5. **Banking & Financial Services**
Calculate settlement dates, interest calculations, and business days excluding holidays.

### 6. **Government & Public Services**
Integrate official holiday data into public service portals and citizen applications.

### 7. **Education & Academic Systems**
Schedule exams, breaks, and academic calendars aligned with national holidays.

---

## 🎯 Who is This For?

| Audience | Use Case |
|----------|----------|
| **Developers** | Build apps with accurate Indonesian holiday data |
| **HR Teams** | Automate leave and payroll calculations |
| **Startups** | Free, reliable holiday API for their products |
| **Enterprise** | Self-hosted solution with full data control |
| **Government** | Public service integration |
| **Students** | Learn Go, REST API, JWT authentication |

---

## 🖥️ Self-Hosting Guide

### Requirements
- Docker & Docker Compose, OR
- Go 1.23+ with SQLite support

### Docker Deployment

#### Production Deployment
```bash
# Clone repository
git clone https://github.com/ilramdhan/holidayapi.git
cd holidayapi

# Create environment file
cp .env.example .env

# Edit .env with your settings
nano .env

# Start production stack
docker-compose -f docker-compose.yml up -d
```

#### Environment Variables
```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Database
DATABASE_PATH=./data/holidays.db
MIGRATIONS_PATH=./migrations

# Security
JWT_SECRET_KEY=your-super-secret-key-min-32-chars
ADMIN_API_KEY=your-admin-api-key

# Rate Limiting
RATE_LIMIT_RPM=60
RATE_LIMIT_BURST=10

# JWT Settings
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=168h
```

### Cloud Deployment

#### Railway
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/your-template-id)

#### Render
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)

#### Fly.io
```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Launch app
fly launch

# Deploy
fly deploy
```

---

## 🔧 SDK & Clients

### Go SDK

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ilramdhan/holidayapi/pkg/client"
)

func main() {
    // Create client
    c := client.New("https://api.holidayapi.id/v1")
    
    // Get holidays for 2024
    holidays, err := c.GetHolidaysByYear(context.Background(), 2024)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, h := range holidays {
        fmt.Printf("%s: %s\n", h.Date, h.Name)
    }
}
```

### JavaScript/TypeScript Client

```typescript
import { HolidayClient } from '@holidayapi/indonesia';

const client = new HolidayClient('https://api.holidayapi.id/v1');

// Get holidays
const holidays = await client.getHolidaysByYear(2024);
console.log(holidays);
```

### Python Client

```python
from holidayapi import HolidayClient

client = HolidayClient('https://api.holidayapi.id/v1')

# Get holidays for 2024
holidays = client.get_holidays_by_year(2024)
for holiday in holidays:
    print(f"{holiday['date']}: {holiday['name']}")
```

### Postman Collection

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" width="128px">](https://god.gw.postman.com/run-collection/your-collection-id)

📥 [Download Postman Collection](./docs/HolidayAPI.postman_collection.json)

---

## 🔗 Integrations

### n8n Integration
Use our n8n node to automate holiday-aware workflows.

### Zapier Integration
Coming soon - Connect Holiday API with 5000+ apps.

### Google Calendar
Sync Indonesian holidays to your Google Calendar.

### Slack Bot
Get holiday notifications in your Slack workspace.

---

## 📖 Usage Examples

### Get Holidays by Year
```bash
curl -X GET "http://localhost:8080/api/v1/holidays/year/2024" \
  -H "Accept: application/json"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "date": "2024-01-01",
      "name": "Tahun Baru 2024",
      "type": "national",
      "description": "Hari Tahun Baru Masehi"
    },
    {
      "id": 2,
      "date": "2024-02-08",
      "name": "Isra Mikraj Nabi Muhammad SAW",
      "type": "national",
      "description": "Kenaikan Nabi Muhammad SAW"
    }
  ],
  "meta": {
    "total": 2,
    "year": 2024
  }
}
```

### Get Today's Holiday
```bash
curl -X GET "http://localhost:8080/api/v1/holidays/today" \
  -H "Accept: application/json"
```

### Get Upcoming Holidays
```bash
curl -X GET "http://localhost:8080/api/v1/holidays/upcoming?limit=5" \
  -H "Accept: application/json"
```

### Filter by Type
```bash
# Only national holidays
curl -X GET "http://localhost:8080/api/v1/holidays/year/2024?type=national"

# Only collective leave days
curl -X GET "http://localhost:8080/api/v1/holidays/year/2024?type=collective_leave"
```

### Authentication Example
```bash
# 1. Login
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123!"
  }'

# 2. Use token for admin endpoints
curl -X POST "http://localhost:8080/api/v1/admin/holidays" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "date": "2024-12-25",
    "name": "Hari Raya Natal",
    "type": "national",
    "description": "Perayaan kelahiran Yesus Kristus"
  }'
```

---

## 🔧 Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `localhost` | Server host |
| `SERVER_PORT` | `8080` | Server port |
| `DATABASE_PATH` | `./data/holidays.db` | SQLite database path |
| `MIGRATIONS_PATH` | `./migrations` | Database migrations path |
| `RATE_LIMIT_RPM` | `60` | Rate limit requests per minute |
| `RATE_LIMIT_BURST` | `10` | Rate limit burst size |
| `JWT_SECRET_KEY` | `your-secret-key` | JWT signing secret key |
| `JWT_ACCESS_TOKEN_TTL` | `15m` | Access token expiration time |
| `JWT_REFRESH_TOKEN_TTL` | `168h` | Refresh token expiration time (7 days) |
| `ADMIN_API_KEY` | `admin-secret-key` | Legacy admin API key |

---

## 🐳 Docker Deployment

### Build Image
```bash
docker build -t holidayapi:latest .
```

### Run Container
```bash
docker run -d \
  --name holidayapi \
  -p 8080:8080 \
  -e JWT_SECRET_KEY=your-secret-key \
  -v $(pwd)/data:/root/data \
  holidayapi:latest
```

### Docker Compose
```yaml
version: '3.8'

services:
  holidayapi:
    build: .
    ports:
      - "8080:8080"
    environment:
      - JWT_SECRET_KEY=${JWT_SECRET_KEY}
      - ADMIN_API_KEY=${ADMIN_API_KEY}
    volumes:
      - ./data:/root/data
    restart: unless-stopped
```

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/services/...
```

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Holiday data based on official SKB 3 Menteri (Ministerial Decree)
- Built with [Gin Web Framework](https://github.com/gin-gonic/gin)
- Database powered by [modernc/sqlite](https://gitlab.com/cznic/sqlite)

---

## 📞 Support

- 📧 Email: ilramdhan@gmail.com
- 🐛 Issues: [GitHub Issues](https://github.com/ilramdhan/holidayapi/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/ilramdhan/holidayapi/discussions)

---

<p align="center">
  Made with ❤️ in Indonesia 🇮🇩
</p>
