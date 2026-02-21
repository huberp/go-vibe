# API Reference

Complete reference for all go-vibe HTTP endpoints.

## Base URL

```
http://localhost:8080
```

All user management endpoints are versioned under `/v1`.

## Authentication

Protected endpoints require a `Bearer` token in the `Authorization` header.

### Get a Token

```bash
curl -s -X POST http://localhost:8080/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "password": "secret123"}'
```

**Response `200 OK`**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Use the Token

```bash
curl -s http://localhost:8080/v1/users \
  -H "Authorization: Bearer <token>"
```

## Endpoints

### `POST /v1/users` — Register User

Creates a new user account. No authentication required.

**Request body**

```json
{
  "name":     "Alice Smith",
  "email":    "alice@example.com",
  "password": "secret123",
  "role":     "user"
}
```

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| `name` | string | ✅ | min 1 char |
| `email` | string | ✅ | valid email, unique |
| `password` | string | ✅ | min 6 chars |
| `role` | string | ❌ | `user` or `admin`, defaults to `user` |

**Response `201 Created`**

```json
{
  "id":         1,
  "name":       "Alice Smith",
  "email":      "alice@example.com",
  "role":       "user",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Error responses**

| Status | Reason |
|--------|--------|
| `400` | Invalid input / validation failure |
| `409` | Email already registered |

---

### `POST /v1/login` — Login

Authenticates a user and returns a signed JWT token.

**Request body**

```json
{
  "email":    "alice@example.com",
  "password": "secret123"
}
```

**Response `200 OK`**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error responses**

| Status | Reason |
|--------|--------|
| `400` | Missing or malformed request body |
| `401` | Invalid email or password |

---

### `GET /v1/users` — List Users

Returns all registered users. **Requires `admin` role.**

**Headers**

```
Authorization: Bearer <admin-token>
```

**Response `200 OK`**

```json
[
  {
    "id":         1,
    "name":       "Alice Smith",
    "email":      "alice@example.com",
    "role":       "admin",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  },
  {
    "id":         2,
    "name":       "Bob Jones",
    "email":      "bob@example.com",
    "role":       "user",
    "created_at": "2024-01-16T09:30:00Z",
    "updated_at": "2024-01-16T09:30:00Z"
  }
]
```

**Error responses**

| Status | Reason |
|--------|--------|
| `401` | Missing or invalid JWT |
| `403` | Role is not `admin` |

---

### `GET /v1/users/:id` — Get User

Returns a single user. Users may fetch their own profile; admins may fetch any user.

**Path parameters**

| Param | Type | Description |
|-------|------|-------------|
| `id` | integer | User ID |

**Example**

```bash
curl -s http://localhost:8080/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Response `200 OK`**

```json
{
  "id":         1,
  "name":       "Alice Smith",
  "email":      "alice@example.com",
  "role":       "admin",
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Error responses**

| Status | Reason |
|--------|--------|
| `401` | Missing or invalid JWT |
| `403` | Attempting to access another user's profile without admin role |
| `404` | User not found |

---

### `PUT /v1/users/:id` — Update User

Updates a user's name, email, or password. Users may update their own record; admins may update any.

**Request body** (all fields optional)

```json
{
  "name":     "Alice Updated",
  "email":    "alice-new@example.com",
  "password": "newpassword456"
}
```

**Response `200 OK`** — returns the updated user object.

**Error responses**

| Status | Reason |
|--------|--------|
| `400` | Validation failure |
| `401` | Missing or invalid JWT |
| `403` | Attempting to update another user without admin role |
| `404` | User not found |

---

### `DELETE /v1/users/:id` — Delete User

Soft-deletes a user. **Requires `admin` role.**

**Example**

```bash
curl -s -X DELETE http://localhost:8080/v1/users/2 \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

**Response `204 No Content`**

**Error responses**

| Status | Reason |
|--------|--------|
| `401` | Missing or invalid JWT |
| `403` | Role is not `admin` |
| `404` | User not found |

---

### `GET /health` — Health Check

Returns server health status. Used by Kubernetes liveness/readiness probes.

**Response `200 OK`**

```json
{
  "status": "healthy"
}
```

---

### `GET /metrics` — Prometheus Metrics

Exposes Prometheus-formatted metrics for scraping.

```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/v1/users",status="200"} 42
...
```

---

### `GET /swagger/*` — Swagger UI

Interactive API documentation powered by Swagger UI. Browse to `http://localhost:8080/swagger/index.html` in your browser.

## Error Format

All error responses use a consistent JSON structure:

```json
{
  "error": "descriptive error message"
}
```

## Complete curl Examples

```bash
BASE="http://localhost:8080"

# --- Setup ---
# Register admin user
curl -s -X POST $BASE/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Admin","email":"admin@example.com","password":"admin123","role":"admin"}' | jq .

# Login → capture token
TOKEN=$(curl -s -X POST $BASE/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' | jq -r .token)

# --- User operations ---
# Create regular user
curl -s -X POST $BASE/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob","email":"bob@example.com","password":"pass123"}' | jq .

# List all users
curl -s $BASE/v1/users -H "Authorization: Bearer $TOKEN" | jq .

# Get user by ID
curl -s $BASE/v1/users/1 -H "Authorization: Bearer $TOKEN" | jq .

# Update user
curl -s -X PUT $BASE/v1/users/2 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Bob Updated"}' | jq .

# Delete user
curl -s -X DELETE $BASE/v1/users/2 \
  -H "Authorization: Bearer $TOKEN"

# Health
curl -s $BASE/health | jq .
```
