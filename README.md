# Holiday API Indonesia 🇮🇩

API untuk mendapatkan informasi hari libur nasional dan cuti bersama Indonesia berdasarkan SKB 3 Menteri.

## ✨ Features

- ✅ **REST API** dengan versioning (v1)
- ✅ **CRUD operations** untuk admin
- ✅ **Filter berdasarkan jenis** holiday (libur nasional, cuti bersama, atau keduanya)
- ✅ **Filter berdasarkan periode** (tahun, bulan, hari)
- ✅ **Rate limiting** (60 requests/minute)
- ✅ **Swagger documentation**
- ✅ **SQLite database** (pure Go, no CGO required)
- ✅ **Comprehensive logging**
- ✅ **Input validation**
- ✅ **Docker support**
- ✅ **Unit & Integration tests**

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git

### Installation & Running

1. **Clone the repository**
```bash
git clone <repository-url>
cd holidayApi
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Run the application**

**Option 1: Using Go directly**
```bash
# Windows
scripts/run.bat

# Linux/Mac
chmod +x scripts/run.sh
./scripts/run.sh
```

**Option 2: Using environment variables**
```bash
export SERVER_HOST=localhost
export SERVER_PORT=8080
export DATABASE_PATH=./data/holidays.db
export MIGRATIONS_PATH=./migrations
export ADMIN_API_KEY=your-secret-key

go run cmd/server/main.go
```

**Option 3: Using Docker**
```bash
docker-compose up --build
```

4. **Access the API**
- **API Base URL**: http://localhost:8080/api/v1
- **Swagger Documentation**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

## 📋 API Endpoints

### Public Endpoints
- `GET /api/v1/holidays` - Get all holidays with filters
- `GET /api/v1/holidays/year/{year}` - Get holidays by year
- `GET /api/v1/holidays/month/{year}/{month}` - Get holidays by month
- `GET /api/v1/holidays/today` - Get today's holiday (if any)
- `GET /api/v1/holidays/this-year` - Get holidays for current year
- `GET /api/v1/holidays/this-month` - Get holidays for current month
- `GET /api/v1/holidays/upcoming` - Get upcoming holidays

### Admin Endpoints (Protected)
- `POST /api/v1/admin/holidays` - Create new holiday
- `GET /api/v1/admin/holidays/{id}` - Get holiday by ID
- `PUT /api/v1/admin/holidays/{id}` - Update holiday
- `DELETE /api/v1/admin/holidays/{id}` - Delete holiday

## 🔧 Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `localhost` | Server host |
| `SERVER_PORT` | `8080` | Server port |
| `DATABASE_PATH` | `./data/holidays.db` | SQLite database path |
| `MIGRATIONS_PATH` | `./migrations` | Database migrations path |
| `RATE_LIMIT_RPM` | `60` | Rate limit requests per minute |
| `RATE_LIMIT_BURST` | `10` | Rate limit burst size |
| `ADMIN_API_KEY` | `admin-secret-key` | Admin API key |

## 📖 Usage Examples

### Get all holidays for 2024
```bash
curl "http://localhost:8080/api/v1/holidays/year/2024"
```

### Get only national holidays for 2024
```bash
curl "http://localhost:8080/api/v1/holidays/year/2024?type=national"
```

### Get today's holiday
```bash
curl "http://localhost:8080/api/v1/holidays/today"
```

### Create a new holiday (Admin)
```bash
curl -X POST "http://localhost:8080/api/v1/admin/holidays" \
  -H "X-API-Key: admin-secret-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hari Libur Khusus",
    "date": "2024-12-31",
    "type": "national",
    "description": "Hari libur khusus akhir tahun"
  }'
```

## 🏗️ Project Structure

```
holidayapi/
├── cmd/
│   └── server/          # Application entrypoints
├── internal/
│   ├── config/          # Configuration
│   ├── database/        # Database setup and migrations
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   └── services/        # Business logic
├── pkg/
│   └── utils/           # Shared utilities
├── api/
│   └── swagger/         # API documentation
├── migrations/          # Database migrations
├── scripts/             # Helper scripts
└── docs/               # Additional documentation
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific test
go test -v ./internal/services/
```

## 🐳 Docker

```bash
# Build and run with Docker Compose
docker-compose up --build

# Build Docker image
docker build -t holidayapi .

# Run Docker container
docker run -p 8080:8080 holidayapi
```

## 📚 Documentation

- **API Documentation**: Available at `/swagger/index.html` when server is running
- **Detailed API Guide**: See [docs/API.md](docs/API.md)

## 🔒 Security

- **Rate Limiting**: 60 requests per minute per IP
- **Admin Authentication**: API key required for admin endpoints
- **Input Validation**: All inputs are validated
- **CORS**: Configured for cross-origin requests

## 🎯 Holiday Types

- **`national`**: Libur Nasional (National Holiday) - Berdasarkan SKB 3 Menteri
- **`collective_leave`**: Cuti Bersama (Collective Leave) - Berdasarkan SKB 3 Menteri

## 📅 Sample Data

API sudah termasuk data hari libur Indonesia untuk tahun 2024-2025 berdasarkan SKB 3 Menteri:

**Libur Nasional 2024:**
- Tahun Baru Masehi (1 Januari)
- Isra Mikraj (8 Februari)
- Tahun Baru Imlek (10 Februari)
- Hari Raya Nyepi (11 Maret)
- Wafat Isa Almasih (29 Maret)
- Hari Raya Idul Fitri (10-11 April)
- Hari Buruh (1 Mei)
- Kenaikan Isa Almasih (9 Mei)
- Hari Raya Waisak (23 Mei)
- Hari Lahir Pancasila (1 Juni)
- Hari Raya Idul Adha (17 Juni)
- Tahun Baru Islam (7 Juli)
- HUT RI ke-79 (17 Agustus)
- Maulid Nabi (16 September)
- Hari Raya Natal (25 Desember)

**Cuti Bersama 2024:**
- Berbagai tanggal cuti bersama sesuai SKB 3 Menteri

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Data hari libur berdasarkan SKB 3 Menteri Republik Indonesia
- Built with Go, Gin, SQLite, dan Swagger
