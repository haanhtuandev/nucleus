# Nucleus API

A modern, production-ready blog platform backend built with Go. Features user authentication, post management, social following system, and full-text search — all secured with rate limiting and XSS protection.

![Go Version](https://img.shields.io/badge/Go-1.25-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Production%20Ready-success)

---

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

```bash
git clone https://github.com/haanhtuandev/nucleus.git
cd go-boilerplate-webservice
docker compose up --build
```

**Done!** API runs on `http://localhost:8080`, PostgreSQL on `localhost:5432`.

### Option 2: Manual Setup

**Prerequisites:** Go 1.25+, PostgreSQL 16+

```bash
# 1. Clone
git clone https://github.com/haanhtuandev/nucleus.git
cd go-boilerplate-webservice

# 2. Configure
cp .env.example .env
# Edit .env with your DB_URL, SECRET (32+ chars), ALLOWED_ORIGINS

# 3. Start PostgreSQL
docker run -d --name blog-postgres \
  -e POSTGRES_USER=bloguser \
  -e POSTGRES_PASSWORD=blogpassword123 \
  -e POSTGRES_DB=blogdb \
  -p 5432:5432 \
  postgres:16-alpine

# 4. Build & Run
go build -o blog-service ./cmd/blog-service
./blog-service
```

---

## 📋 Features

| Category | Features |
|----------|----------|
| **Authentication** | JWT access + refresh tokens, Argon2id password hashing |
| **User Management** | Registration, login, profile management |
| **Posts** | Create, update, delete, soft deletes, slug-based URLs |
| **Social** | Follow/unfollow users, feed system |
| **Search** | Full-text PostgreSQL search |
| **Security** | Rate limiting (1000 req/hr), XSS protection, CORS, SQL injection prevention |
| **Database** | PostgreSQL 16+, sqlc for type-safe queries, auto-migrations |

---

## 📖 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/signup` | Register new user | ❌ |
| `POST` | `/login` | Login with credentials | ❌ |
| `POST` | `/auth/refresh` | Refresh access token | ❌ |

### Posts

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/me/posts` | Create post | ✅ |
| `PUT` | `/me/posts/:id` | Update post | ✅ |
| `DELETE` | `/me/posts/:id` | Delete post | ✅ |
| `GET` | `/posts` | List all posts | ❌ |
| `GET` | `/posts/:id` | Get post by ID | ❌ |
| `GET` | `/posts/search` | Search posts | ❌ |

### Users & Social

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/me/profile` | Get my profile | ✅ |
| `GET` | `/profile/:id` | Get user profile | ❌ |
| `GET` | `/follow/:id` | Follow user | ✅ |
| `GET` | `/unfollow/:id` | Unfollow user | ✅ |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |

---

## 🧪 Testing the API

```bash
# Health check
curl http://localhost:8080/health

# Sign up
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","password":"securepass123","bio":"Developer"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","password":"securepass123"}'

# Create post (use token from login)
curl -X POST http://localhost:8080/me/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"My Post","content":"Post content here"}'

# Search posts
curl "http://localhost:8080/posts/search?q=go"
```

See `test.http` for more examples.

---

## 🛠️ Tech Stack

| Component | Technology |
|-----------|------------|
| **Language** | Go 1.25 |
| **Database** | PostgreSQL 16 |
| **ORM** | sqlc (type-safe SQL generation) |
| **Migrations** | goose (embedded) |
| **Auth** | JWT (HS256) + refresh tokens |
| **Password Hashing** | Argon2id |
| **Containerization** | Docker + Docker Compose |

---

## 📁 Project Structure

```
go-boilerplate-webservice/
├── cmd/
│   └── blog-service/       # Application entry point
├── internal/
│   ├── api/                # HTTP handlers, routes, middleware
│   ├── auth/               # JWT, password hashing (Argon2id)
│   └── database/           # sqlc generated code, migrations
├── sql/
│   ├── schema/             # Database migrations (001_*.sql, etc.)
│   └── queries/            # SQL queries for sqlc generation
├── static/                 # Static assets (homepage)
├── .env.example            # Environment variable template
├── docker-compose.yml      # Docker orchestration
├── Dockerfile              # Container build instructions
└── sqlc.yaml               # sqlc configuration
```

---

## ⚙️ Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Database connection
DB_URL=postgres://bloguser:blogpassword123@localhost:5432/blogdb?sslmode=disable

# JWT secret (MUST be 32+ characters)
# Generate with: openssl rand -base64 32
SECRET=your-super-secret-key-at-least-32-characters-long

# Allowed CORS origins (comma-separated)
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Docker Configuration

Override defaults via environment variables in `docker-compose.yml`:

```bash
# Custom secret and CORS
SECRET=your-secret ALLOWED_ORIGINS=https://myapp.com docker compose up
```

---

## 🚢 Deployment

### Deploy to Render (Free Tier)

1. **Push to GitHub**
   ```bash
   git push origin main
   ```

2. **Create PostgreSQL on Render**
   - New → PostgreSQL → Free tier
   - Copy **Internal Database URL**

3. **Create Web Service**
   - New → Web Service → Connect GitHub repo
   - **Build Command:** `go build -o blog-service ./cmd/blog-service`
   - **Start Command:** `./blog-service`
   - **Environment Variables:**
     ```
     DB_URL=<paste Internal Database URL>
     SECRET=<generate: openssl rand -base64 32>
     ALLOWED_ORIGINS=https://your-app.onrender.com
     ```

4. **Deploy!** Your API is live at `https://your-app.onrender.com`

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed guides on:
- Docker Compose production deployment
- VPS deployment (DigitalOcean, Linode)
- Google Cloud Run
- AWS deployment

---

## 🔒 Security Features

| Feature | Implementation |
|---------|----------------|
| **Password Storage** | Argon2id with memory-hard parameters |
| **JWT Validation** | Algorithm validation (prevents alg:none attacks) |
| **Rate Limiting** | Sliding window, per-IP tracking (1000 req/hr) |
| **XSS Prevention** | HTML tag stripping on user input |
| **SQL Injection** | Parameterized queries via sqlc |
| **CORS** | Configurable allowed origins |
| **Soft Deletes** | Posts are soft-deleted (recoverable) |

---

## 🧪 Development

### Run Tests

```bash
go test ./...
```

### Generate SQLC Code

After modifying `sql/queries/*.sql`:

```bash
sqlc generate
```

### Create New Migration

```bash
# Create new migration file
touch sql/schema/010_your_migration_name.sql

# Add Up and Down sections:
# -- +goose Up
# CREATE TABLE ...
#
# -- +goose Down
# DROP TABLE ...
```

Migrations run automatically on application startup.

---

## 📄 License

MIT License - see LICENSE file for details.

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/your-feature`)
3. Commit changes (`git commit -m 'Add your feature'`)
4. Push to branch (`git push origin feature/your-feature`)
5. Open Pull Request

---

## 📞 Support

- **Issues:** [GitHub Issues](https://github.com/yourusername/go-boilerplate-webservice/issues)
- **Deployment Guide:** [DEPLOYMENT.md](DEPLOYMENT.md)
- **API Examples:** [test.http](test.http)
