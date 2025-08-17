# Holiday API Indonesia - API Documentation

API untuk mendapatkan informasi hari libur nasional dan cuti bersama Indonesia berdasarkan SKB 3 Menteri dengan sistem authentication JWT yang lengkap.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

### JWT Authentication (Recommended)
API menggunakan JWT (JSON Web Token) untuk authentication. Setelah login, Anda akan mendapatkan access token dan refresh token.

**Authorization Header:**
```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

### Default Admin Credentials
```
Username: admin
Password: Admin123!
Role: super_admin
```

### Legacy API Key (Backward Compatibility)
Admin endpoints juga masih mendukung API Key untuk backward compatibility:
```
X-API-Key: your-admin-api-key
```

## Rate Limiting

- **Public endpoints**: 60 requests per minute
- **Burst limit**: 10 requests
- Rate limit berdasarkan IP address

## Endpoints

### Authentication Endpoints

#### Login
```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "username": "admin",
  "password": "Admin123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@holidayapi.com",
      "role": "super_admin"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### Refresh Token
```http
POST /api/v1/auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "your-refresh-token"
}
```

#### Get Profile (JWT Required)
```http
GET /api/v1/auth/profile
```

#### Change Password (JWT Required)
```http
POST /api/v1/auth/change-password
```

**Request Body:**
```json
{
  "current_password": "Admin123!",
  "new_password": "NewPassword123!"
}
```

#### Register New User (Super Admin Only)
```http
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "username": "newadmin",
  "email": "newadmin@holidayapi.com",
  "password": "NewAdmin123!",
  "role": "admin"
}
```

### Public Endpoints (No Authentication Required)

#### 1. Get All Holidays
```http
GET /api/v1/holidays
```

**Query Parameters:**
- `year` (int, optional): Filter by year
- `month` (int, optional): Filter by month (1-12)
- `type` (string, optional): Filter by type (`national` or `collective_leave`)
- `limit` (int, optional): Limit results (default: 50, max: 100)
- `offset` (int, optional): Offset for pagination (default: 0)

**Example:**
```http
GET /api/v1/holidays?year=2024&type=national&limit=10
```

#### 2. Get Holidays by Year
```http
GET /api/v1/holidays/year/{year}
```

**Parameters:**
- `year` (int): Year to filter

**Query Parameters:**
- `type` (string, optional): Filter by type

**Example:**
```http
GET /api/v1/holidays/year/2024?type=collective_leave
```

#### 3. Get Holidays by Month
```http
GET /api/v1/holidays/month/{year}/{month}
```

**Parameters:**
- `year` (int): Year
- `month` (int): Month (1-12)

**Query Parameters:**
- `type` (string, optional): Filter by type

#### 4. Get Today's Holiday
```http
GET /api/v1/holidays/today
```

Returns holiday for today's date if any exists.

#### 5. Get This Year's Holidays
```http
GET /api/v1/holidays/this-year
```

**Query Parameters:**
- `type` (string, optional): Filter by type

#### 6. Get This Month's Holidays
```http
GET /api/v1/holidays/this-month
```

**Query Parameters:**
- `type` (string, optional): Filter by type

#### 7. Get Upcoming Holidays
```http
GET /api/v1/holidays/upcoming
```

**Query Parameters:**
- `limit` (int, optional): Limit results (default: 10)
- `type` (string, optional): Filter by type

### Admin Endpoints (JWT Required - Admin/Super Admin)

#### 1. Create Holiday
```http
POST /api/v1/admin/holidays
```

**Headers:**
```
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Hari Raya Idul Fitri",
  "date": "2024-04-10",
  "type": "national",
  "description": "Hari libur nasional Hari Raya Idul Fitri 1445 Hijriah"
}
```

#### 2. Update Holiday
```http
PUT /api/v1/admin/holidays/{id}
```

**Headers:**
```
Authorization: Bearer YOUR_ACCESS_TOKEN
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Updated Holiday Name",
  "description": "Updated description",
  "is_active": true
}
```

#### 3. Delete Holiday
```http
DELETE /api/v1/admin/holidays/{id}
```

**Headers:**
```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

#### 4. Get Audit Logs (Admin Only)
```http
GET /api/v1/admin/audit-logs
```

**Headers:**
```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

**Query Parameters:**
- `user_id` (int, optional): Filter by user ID
- `action` (string, optional): Filter by action type
- `resource` (string, optional): Filter by resource type
- `success` (bool, optional): Filter by success status
- `start_date` (string, optional): Start date (YYYY-MM-DD)
- `end_date` (string, optional): End date (YYYY-MM-DD)
- `limit` (int, optional): Limit results (default: 50, max: 100)
- `offset` (int, optional): Offset for pagination (default: 0)

#### 5. Get User Audit Logs (Admin Only)
```http
GET /api/v1/admin/audit-logs/user/{id}
```

**Headers:**
```
Authorization: Bearer YOUR_ACCESS_TOKEN
```

## Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error description"
}
```

### Paginated Response
```json
{
  "success": true,
  "message": "Data retrieved successfully",
  "data": {
    "data": [
      // Array of holidays
    ],
    "total": 25,
    "page": 1,
    "per_page": 10,
    "total_pages": 3
  }
}
```

## Holiday Types

- `national`: Libur Nasional (National Holiday)
- `collective_leave`: Cuti Bersama (Collective Leave)

## Date Format

All dates use ISO 8601 format: `YYYY-MM-DD`

## Error Codes

- `400`: Bad Request - Invalid input
- `401`: Unauthorized - Missing or invalid API key
- `404`: Not Found - Resource not found
- `429`: Too Many Requests - Rate limit exceeded
- `500`: Internal Server Error - Server error

## Examples

### Get all national holidays for 2024
```bash
curl "http://localhost:8080/api/v1/holidays?year=2024&type=national"
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
