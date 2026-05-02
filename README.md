# Balloon: A complete backend boilerplate written in Go

A production-ready monorepo template for building scalable web applications with Go backend and TypeScript frontend. Built with modern best practices, clean architecture, and comprehensive tooling.

# Features
Uses Echo as the main HTTP framework, wiring clean routes,
middlewares, and handler layers for REST APIs.

Integrated PostgreSQL and Redis for data storage and caching
with reusable modules for CRUD, sessions, and background
jobs, while defining OpenAPI specifications for core services to
enable automatic documentation and easier client generation.

Added Relic for monitoring and metrics, implemented Clerk-
based authentication for sign-up, sign-in, and sessions, and set
up Testcontainers-based tests that spin up real Postgres/Redis
instances in CI to improve integration test reliability.

# Project Structure
```
go-boilerplate/
├── apps/backend/     # All the Go backend application
├── packages/         # Frontend packages (Typescript, React, etc)
├── package.json      # Monorepo configuration
├── turbo.json        # Turborepo configuration
└── README.md         # Current file
```


# Quick Start 
There are some prerequisites to start up the project.
- Go 1.24 or higher
- Node.js 22+ and Bun
- PostgreSQL 16+
- Redis 8+

**Installation process**:-

1. **Clone this repository**:
```
git clone https://github.com/nrngxv/balloon.git
cd balloon
```

2. **Install dependencies**:
```
# Install frontend dependencies
bun install

# Install backend dependencies
cd apps/backend
go mod download
```

3. Set up environment variables:
```
cp apps/backend/.env.example apps/backend/.env
# Edit apps/backend/.env with your configuration
```

4. Start database and Redis:

5. Run the database migration:
```
cd apps/backend
task migrations:up
```

6. Start development server:
```
# From root directory
bun dev

# Or just the backend
cd apps/backend
task run
```

The API will be available at `http://localhost:8080`

# Development

## Available commands
```
# Backend commands (from backend/ directory)
task help              # Show all available tasks
task run               # Run the application
task migrations:new    # Create a new migration
task migrations:up     # Apply migrations
task test              # Run tests
task tidy              # Format code and manage dependencies

# Frontend commands (from root directory)
bun dev                # Start development servers
bun build              # Build all packages
bun lint               # Lint all packages
```

## Environment variables
The backend uses environment variables prefixed with `BALLOON_`. Key variables include:

- `BALLOON_DATABASE_*` - PostgreSQL connection settings
- `BALLOON_SERVER_*` - Server configuration
- `BALLOON_AUTH_*` - Authentication settings
- `BALLOON_REDIS_*` - Redis connection
- `BALLOON_EMAIL_*` - Email service configuration
- `BALLOON_OBSERVABILITY_*` - Monitoring settings

See `apps/backend/.env.example` for a complete list.

# Architecture
This boilerplate/template follows clean architecture principles:

- **Handlers**: HTTP request/response handling
- **Services**: Business logic implementation
- **Repositories**: Data access layer
- **Models**: Domain entities
- **Infrastructure**: External services (Database, cache, email)

# Testing

```
# Run backend tests
cd apps/backend
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (requires Docker)
go test -tags=integration ./...
```

### Production Considerations
1. Use environment-specific configuration
2. Enable production logging levels
3. Configure proper database connection pooling
4. Set up monitoring and alerting
5. Use a reverse proxy (nginx, Caddy)
6. Enable rate limiting and security headers
7. Configure CORS for your domains
