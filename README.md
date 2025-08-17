# Holiday API Indonesia 🇮🇩

API untuk mendapatkan informasi hari libur nasional dan cuti bersama Indonesia berdasarkan SKB 3 Menteri.

## ✨ Features

### 🔐 **Authentication & Security**
- ✅ **JWT Authentication** dengan access & refresh tokens
- ✅ **Role-Based Access Control (RBAC)** - Super Admin & Admin roles
- ✅ **User Management** - Registration, login, profile management
- ✅ **Password Security** - Bcrypt hashing dengan password policy
- ✅ **Comprehensive Audit Logging** - Track semua user actions
- ✅ **Security Headers** - XSS, CSRF, Content Security Policy protection
- ✅ **Input Sanitization** - Protection dari injection attacks
- ✅ **Enhanced Rate Limiting** - Per-user rate limiting

### 🚀 **API Features**
- ✅ **REST API** dengan versioning (v1)
- ✅ **CRUD operations** untuk admin (JWT protected)
- ✅ **Filter berdasarkan jenis** holiday (libur nasional, cuti bersama, atau keduanya)
- ✅ **Filter berdasarkan periode** (tahun, bulan, hari)
- ✅ **Swagger documentation** dengan authentication
- ✅ **SQLite database** (pure Go, no CGO required)
- ✅ **Comprehensive logging** dengan structured format
- ✅ **Input validation** dengan custom validators
- ✅ **Docker support** dengan production-ready configuration
- ✅ **Unit & Integration tests** dengan mocking

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

5. **Setup Admin User**
```bash
# First, create your admin user via API or database
# See DEPLOYMENT.md for secure setup instructions
```

6. **Login & Get JWT Token**
```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "your-admin-username",
    "password": "your-secure-password"
  }'
```

## 📋 API Endpoints

### 🔓 **Public Endpoints (No Authentication Required)**
```
GET /api/v1/holidays                    - Get all holidays with filters
GET /api/v1/holidays/year/{year}        - Get holidays by year
GET /api/v1/holidays/month/{year}/{month} - Get holidays by month
GET /api/v1/holidays/today              - Get today's holiday (if any)
GET /api/v1/holidays/this-year          - Get holidays for current year
GET /api/v1/holidays/this-month         - Get holidays for current month
GET /api/v1/holidays/upcoming           - Get upcoming holidays
GET /health                             - Health check endpoint
```

### 🔐 **Authentication Endpoints**
```
POST /api/v1/auth/login                 - User login (get JWT tokens)
POST /api/v1/auth/refresh               - Refresh access token
GET  /api/v1/auth/profile               - Get user profile (JWT required)
POST /api/v1/auth/change-password       - Change password (JWT required)
GET  /api/v1/auth/audit-logs            - Get my audit logs (JWT required)
```

### 👑 **Super Admin Only Endpoints**
```
POST /api/v1/auth/register              - Register new user
GET  /api/v1/auth/users                 - Get all users
DELETE /api/v1/auth/users/{id}          - Delete user
```

### 🛡️ **Admin Endpoints (JWT Required - Admin/Super Admin)**
```
POST /api/v1/admin/holidays             - Create new holiday
GET  /api/v1/admin/holidays/{id}        - Get holiday by ID
PUT  /api/v1/admin/holidays/{id}        - Update holiday
DELETE /api/v1/admin/holidays/{id}      - Delete holiday
GET  /api/v1/admin/audit-logs           - Get all audit logs
GET  /api/v1/admin/audit-logs/user/{id} - Get user audit logs
```

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
| `JWT_SECRET_KEY` | `your-secret-key` | JWT signing secret key |
| `JWT_ACCESS_TOKEN_TTL` | `15m` | Access token expiration time |
| `JWT_REFRESH_TOKEN_TTL` | `168h` | Refresh token expiration time (7 days) |
| `ADMIN_API_KEY` | `admin-secret-key` | Legacy admin API key (backward compatibility) |

## 📖 Usage Examples

### 🔓 **Public API Usage (No Authentication)**

#### Get all holidays for 2024
```bash
curl "http://localhost:8080/api/v1/holidays/year/2024"
```

#### Get only national holidays for 2024
```bash
curl "http://localhost:8080/api/v1/holidays/year/2024?type=national"
```

#### Get today's holiday
```bash
curl "http://localhost:8080/api/v1/holidays/today"
```

### 🔐 **Authentication Examples**

#### 1. Login and get JWT tokens
```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123!"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "role": "super_admin"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### 2. Use JWT token for admin operations
```bash
# Save the access_token from login response
ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Create a new holiday
curl -X POST "http://localhost:8080/api/v1/admin/holidays" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hari Libur Khusus",
    "date": "2024-12-31",
    "type": "national",
    "description": "Hari libur khusus akhir tahun"
  }'
```

#### 3. Get user profile
```bash
curl -X GET "http://localhost:8080/api/v1/auth/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

#### 4. View audit logs
```bash
curl -X GET "http://localhost:8080/api/v1/admin/audit-logs" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

#### 5. Register new user (Super Admin only)
```bash
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newadmin",
    "email": "newadmin@holidayapi.com",
    "password": "NewAdmin123!",
    "role": "admin"
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

## 🔒 Security Features

### 🛡️ **Authentication & Authorization**
- **JWT Authentication**: Secure token-based authentication
- **Role-Based Access Control**: Super Admin, Admin roles with different permissions
- **Password Security**: Bcrypt hashing with strong password policy
- **Token Management**: Access tokens (15min) + Refresh tokens (7 days)
- **Session Security**: Secure token validation and refresh mechanism

### 🔐 **Security Middleware**
- **Security Headers**: XSS protection, CSRF protection, Content Security Policy
- **Input Sanitization**: Protection from XSS and injection attacks
- **SQL Injection Protection**: Pattern detection and prevention
- **Rate Limiting**: 60 requests per minute per user/IP
- **CORS Protection**: Proper cross-origin resource sharing
- **Request Validation**: Comprehensive input validation

### 📊 **Audit & Monitoring**
- **Comprehensive Audit Logging**: All user actions tracked
- **Security Events**: Login attempts, failed authentications
- **User Activity**: CRUD operations, system access
- **Audit Trail**: User, IP, User-Agent, timestamp, success/failure status

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
