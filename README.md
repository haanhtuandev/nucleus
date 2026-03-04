# Nucleus API

A modern, production-ready blog platform backend built with Go. Features user authentication, post management, social following system, and full-text search — all secured with rate limiting and XSS protection.

![Go Version](https://img.shields.io/badge/Go-1.24-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Production%20Ready-success)

---

## Description

**Nucleus API** is a RESTful backend service for building blog platforms. It handles everything from user registration and authentication to creating posts, following other writers, and discovering content through search.



Built with performance and security in mind, it includes:
- JWT-based authentication with refresh tokens
- PostgreSQL database with sqlc for type-safe queries
- Rate limiting to prevent abuse
- XSS protection through input sanitization
- CORS support for frontend integration

---

## Motivation

I built Nucleus API to practice backend development fundamentals while creating something practical and production-ready. This project showcases:

- **Moderately Clean Architecture** — Separation of concerns with handlers, services, and database layers
- **Security Best Practices** — Password hashing (Argon2id), JWT validation, rate limiting, input sanitization
- **Database Design** — Proper indexing, foreign keys, soft deletes, full-text search
- **API Design** — RESTful endpoints, consistent error handling, pagination

Without using frameworks, I had to implement many features manually, which gives me a good understanding of a lower-level backend system.

---

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- [goose](https://github.com/pressly/goose) (for migrations)

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/go-boilerplate-webservice.git
cd go-boilerplate-webservice
```

### 2. Set Up Environment Variables

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
DB_URL=postgres://user:password@localhost:5432/blogdb?sslmode=disable
SECRET=your-super-secret-key-at-least-32-characters-long
ALLOWED_ORIGINS=http://localhost:3000
```

### 3. Start PostgreSQL

```bash
# Using Docker
docker run -d --name blog-postgres \
  -e POSTGRES_USER=bloguser \
  -e POSTGRES_PASSWORD=blogpassword123 \
  -e POSTGRES_DB=blogdb \
  -p 5432:5432 \
  postgres:16-alpine
```

### 4. Run Database Migrations

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run all migrations
goose -dir sql/schema postgres "postgres://bloguser:blogpassword123@localhost:5432/blogdb?sslmode=disable" up
```

### 5. Build and Run

```bash
# Build
go build -o blog-service ./cmd/blog-service

# Run
./blog-service
```

The API server starts on `http://localhost:8080`

### 6. Test It Works

```bash
curl http://localhost:8080/health
# Response: Service is healthy!
```

---

## Usage

### Authentication Flow

#### 1. Sign Up

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123",
    "bio": "Software developer from NYC"
  }'
```

#### 2. Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

**Response:**
```json
{
  "id": "d050cbf9-1561-4b34-b473-17664d87838c",
  "username": "johndoe",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "4267596a25ad90f72b5977229d9c7330..."
}
```

#### 3. Create a Post (Authenticated)

```bash
curl -X POST http://localhost:8080/me/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "How I Learned Go",
    "content": "Today I want to share my journey..."
  }'
```

### Key Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `POST` | `/signup` | Register new user | ❌ |
| `POST` | `/login` | Authenticate user | ❌ |
| `POST` | `/auth/refresh` | Refresh access token | ❌ |
| `POST` | `/me/posts` | Create a post | ✅ |
| `PUT` | `/me/posts/:id` | Update a post | ✅ |
| `DELETE` | `/me/posts/:id` | Delete a post | ✅ |
| `GET` | `/posts` | Get all posts | ❌ |
| `GET` | `/posts/:slug` | Get post by slug | ❌ |
| `GET` | `/posts/search` | Search posts | ❌ |
| `GET` | `/me/profile` | Get my profile | ✅ |
| `GET` | `/profile/:id` | Get user profile | ❌ |
| `GET` | `/follow/:id` | Follow a user | ✅ |
| `GET` | `/unfollow/:id` | Unfollow a user | ✅ |

### Rate Limiting

All endpoints are rate-limited to **1,000 requests per hour per IP** to prevent abuse.

---

## 🛠️ Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.24 |
| **Database** | PostgreSQL 16 |
| **ORM** | sqlc (type-safe SQL) |
| **Authentication** | JWT (HS256) + Refresh Tokens |
| **Password Hashing** | Argon2id |
| **Migrations** | goose |
| **Containerization** | Docker |
| **Deployment** | Render / Railway / GCP Cloud Run |

---

## 📁 Project Structure

```
go-boilerplate-webservice/
├── cmd/
│   └── blog-service/       # Application entry point
├── internal/
│   ├── api/                # HTTP handlers, middleware, routes
│   ├── auth/               # JWT, password hashing
│   ├── database/           # sqlc generated code
│   └── middleware/         # CORS, rate limiting
├── sql/
│   ├── schema/             # Database migrations
│   └── queries/            # SQL queries for sqlc
├── static/                 # Static files (homepage)
├── Dockerfile              # Container configuration
├── .env.example            # Environment template
└── swagger.yaml            # API documentation
```

---

## Security Features

| Feature | Implementation |
|---------|----------------|
| **Password Storage** | Argon2id with memory-hard parameters |
| **JWT Validation** | Algorithm validation (prevents alg:none attacks) |
| **Rate Limiting** | Sliding window, per-IP tracking |
| **XSS Prevention** | HTML tag stripping on user input |
| **SQL Injection** | Parameterized queries via sqlc |
| **CORS** | Configurable allowed origins |
| **Soft Deletes** | Posts are soft-deleted (recoverable) |


---


### Docker

```bash
# Build image
docker build -t blog-service .

# Run container
docker run -p 8080:8080 --env-file .env blog-service
```

## Contributing

Contributions are welcome! Here's how you can help:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Guidelines

- Follow existing code style
- Add tests for new features
- Update documentation as needed
- Keep commits atomic and well-described



