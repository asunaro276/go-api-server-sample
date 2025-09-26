# API Documentation

## Overview

This document provides comprehensive documentation for the Go API Server with User Management. The API follows REST principles and provides CRUD operations for user management.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, no authentication is required for API endpoints. This is suitable for development and demonstration purposes.

## Common Response Format

### Success Response

All successful API responses follow a consistent format:

```json
{
  "id": 1,
  "name": "田中太郎",
  "email": "tanaka@example.com",
  "created_at": "2023-12-01T12:00:00Z",
  "updated_at": "2023-12-01T12:00:00Z"
}
```

### Error Response

All error responses follow a consistent format:

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message",
  "details": [
    {
      "field": "email",
      "message": "Field-specific error message"
    }
  ]
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `USER_NOT_FOUND` | Requested user does not exist |
| `EMAIL_ALREADY_EXISTS` | Email address is already in use |
| `INTERNAL_ERROR` | Internal server error |

## Endpoints

### Health Check

#### GET /health

Check the health status of the API server.

**Response:**
- **200 OK**: Server is healthy

```json
{
  "status": "ok",
  "timestamp": "2023-12-01T12:00:00Z"
}
```

### User Management

#### POST /api/v1/users

Create a new user.

**Request Body:**
```json
{
  "name": "田中太郎",
  "email": "tanaka@example.com"
}
```

**Validation Rules:**
- `name`: Required, maximum 100 characters
- `email`: Required, valid email format, maximum 255 characters, must be unique

**Responses:**
- **201 Created**: User created successfully
- **400 Bad Request**: Validation error
- **409 Conflict**: Email already exists
- **500 Internal Server Error**: Server error

**Example Success Response:**
```json
{
  "id": 1,
  "name": "田中太郎",
  "email": "tanaka@example.com",
  "created_at": "2023-12-01T12:00:00Z",
  "updated_at": "2023-12-01T12:00:00Z"
}
```

#### GET /api/v1/users

Retrieve a paginated list of users.

**Query Parameters:**
- `limit` (optional): Number of users to retrieve (default: 10, max: 100)
- `offset` (optional): Number of users to skip (default: 0)

**Example Request:**
```
GET /api/v1/users?limit=5&offset=0
```

**Responses:**
- **200 OK**: Users retrieved successfully
- **500 Internal Server Error**: Server error

**Example Success Response:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "田中太郎",
      "email": "tanaka@example.com",
      "created_at": "2023-12-01T12:00:00Z",
      "updated_at": "2023-12-01T12:00:00Z"
    },
    {
      "id": 2,
      "name": "佐藤花子",
      "email": "sato@example.com",
      "created_at": "2023-12-01T12:30:00Z",
      "updated_at": "2023-12-01T12:30:00Z"
    }
  ],
  "total": 100,
  "limit": 5,
  "offset": 0
}
```

#### GET /api/v1/users/{id}

Retrieve a specific user by ID.

**Path Parameters:**
- `id`: User ID (positive integer)

**Responses:**
- **200 OK**: User retrieved successfully
- **400 Bad Request**: Invalid user ID
- **404 Not Found**: User not found
- **500 Internal Server Error**: Server error

**Example Success Response:**
```json
{
  "id": 1,
  "name": "田中太郎",
  "email": "tanaka@example.com",
  "created_at": "2023-12-01T12:00:00Z",
  "updated_at": "2023-12-01T12:00:00Z"
}
```

#### PUT /api/v1/users/{id}

Update an existing user.

**Path Parameters:**
- `id`: User ID (positive integer)

**Request Body:**
```json
{
  "name": "田中次郎",
  "email": "tanaka.jiro@example.com"
}
```

**Validation Rules:**
- `name` (optional): Maximum 100 characters
- `email` (optional): Valid email format, maximum 255 characters, must be unique

**Responses:**
- **200 OK**: User updated successfully
- **400 Bad Request**: Validation error or invalid user ID
- **404 Not Found**: User not found
- **409 Conflict**: Email already exists
- **500 Internal Server Error**: Server error

**Example Success Response:**
```json
{
  "id": 1,
  "name": "田中次郎",
  "email": "tanaka.jiro@example.com",
  "created_at": "2023-12-01T12:00:00Z",
  "updated_at": "2023-12-01T12:45:00Z"
}
```

#### DELETE /api/v1/users/{id}

Delete a user (soft delete).

**Path Parameters:**
- `id`: User ID (positive integer)

**Responses:**
- **204 No Content**: User deleted successfully
- **400 Bad Request**: Invalid user ID
- **404 Not Found**: User not found
- **500 Internal Server Error**: Server error

## Rate Limiting

Currently, no rate limiting is implemented. This should be added for production use.

## Pagination

List endpoints support pagination using `limit` and `offset` parameters:

- `limit`: Controls the number of items returned (1-100, default: 10)
- `offset`: Controls the number of items to skip (default: 0)

The response includes pagination metadata:

```json
{
  "users": [...],
  "total": 150,
  "limit": 10,
  "offset": 20
}
```

## Data Validation

### User Entity Validation

- **Name**: Required, 1-100 characters, whitespace trimmed
- **Email**: Required, valid email format, case-insensitive, 1-255 characters, must be unique

### Error Handling

Validation errors return detailed information about failed fields:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "リクエストのバリデーションに失敗しました",
  "details": [
    {
      "field": "email",
      "message": "有効なメールアドレスを入力してください"
    }
  ]
}
```

## Examples

### Create User Flow

1. **Create User**
   ```bash
   curl -X POST http://localhost:8080/api/v1/users \
     -H "Content-Type: application/json" \
     -d '{"name":"田中太郎","email":"tanaka@example.com"}'
   ```

2. **Response**
   ```json
   {
     "id": 1,
     "name": "田中太郎",
     "email": "tanaka@example.com",
     "created_at": "2023-12-01T12:00:00Z",
     "updated_at": "2023-12-01T12:00:00Z"
   }
   ```

### Complete CRUD Flow

```bash
# Create
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"田中太郎","email":"tanaka@example.com"}'

# Read (by ID)
curl http://localhost:8080/api/v1/users/1

# Read (list)
curl http://localhost:8080/api/v1/users?limit=10&offset=0

# Update
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"田中次郎","email":"tanaka.jiro@example.com"}'

# Delete
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## OpenAPI Specification

The complete OpenAPI 3.0 specification is available at:
- File: `/specs/001-go-api-api/contracts/user-api.yaml`
- Future endpoint (not implemented): `GET /api/v1/docs`

## Testing

### Manual Testing

Use the examples above or tools like Postman, Insomnia, or curl to test the API endpoints.

### Automated Testing

The project includes comprehensive test suites:

- **Unit Tests**: Test individual components
- **Integration Tests**: Test API endpoints end-to-end
- **Contract Tests**: Validate API contracts
- **Performance Tests**: Ensure response times < 200ms

Run tests:
```bash
go test ./...
```

Run with coverage:
```bash
./scripts/coverage.sh
```

## Performance

### Response Time Targets

All API endpoints are designed to respond within:
- **Create User**: < 200ms
- **Get User**: < 200ms
- **Update User**: < 200ms
- **Delete User**: < 200ms
- **List Users**: < 200ms

### Concurrency

The API is designed to handle concurrent requests efficiently with proper database connection pooling and goroutine management.

## Monitoring

### Health Checks

- **Basic Health**: `GET /health`
- **Detailed Health**: `GET /health/detailed` (future)
- **Readiness**: `GET /health/ready` (future)
- **Liveness**: `GET /health/live` (future)

### Metrics

Future implementation will include:
- Request rate and response times
- Error rates by endpoint
- Database connection pool metrics
- Memory and CPU usage

## Security Considerations

### Current State (Development)

- No authentication required
- No rate limiting
- Basic input validation
- CORS enabled for all origins

### Production Recommendations

- Implement authentication (JWT/OAuth2)
- Add rate limiting
- Restrict CORS origins
- Use HTTPS only
- Add request logging and monitoring
- Implement API versioning strategy
- Add input sanitization

## Migration and Deployment

### Database Migrations

```bash
# Run migrations
go run cmd/api-server/main.go migrate

# Reset database (development only)
go run cmd/api-server/main.go migrate-reset
```

### Environment Configuration

Required environment variables:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `SERVER_HOST`, `SERVER_PORT`
- `APP_ENV` (development/production)

See `quickstart.md` for complete setup instructions.