# Holiday API Indonesia - Documentation

API untuk mendapatkan informasi hari libur nasional dan cuti bersama Indonesia berdasarkan SKB 3 Menteri.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

Admin endpoints memerlukan API Key yang dikirim melalui header:
```
X-API-Key: your-admin-api-key
```

## Rate Limiting

- **Public endpoints**: 60 requests per minute
- **Burst limit**: 10 requests
- Rate limit berdasarkan IP address

## Endpoints

### Public Endpoints

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

### Admin Endpoints (Protected)

#### 1. Create Holiday
```http
POST /api/v1/admin/holidays
```

**Headers:**
```
X-API-Key: your-admin-api-key
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
X-API-Key: your-admin-api-key
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
X-API-Key: your-admin-api-key
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
