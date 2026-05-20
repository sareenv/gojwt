# Go JWT Authentication Service

A professional JWT-based authentication service built with Go, Gin, and PostgreSQL. This project implements secure user registration, login, and token management (access & refresh tokens) with modern observability features.

## Features

- **JWT Authentication**: Secure user registration, login, and logout.
- **Token Management**: Issuance and revocation of access and refresh tokens.
- **Structured Logging**: JSON-based logging using Go's modern `log/slog` library.
- **OpenAPI 3.0**: Comprehensive API documentation available in `openapi.yaml`.
- **Database Integration**: High-performance PostgreSQL integration using `pgx/v5`.
- **Environment Configuration**: Flexible configuration via environment variables or `.env` file.

## Prerequisites

- **Go**: 1.25 or higher
- **PostgreSQL**: 15 or higher
- **Docker & Docker Compose**: For containerized deployment and observability stack.

## Getting Started

### 1. Clone the repository
```bash
git clone <repository-url>
cd gojwt
```

### 2. Configure Environment Variables
Create a `.env` file in the root directory (refer to the existing `.env` or the list below):
- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
- `JWT_SECRET`, `ACCESS_TOKEN_DURATION`, `REFRESH_TOKEN_DURATION`

### 3. Run with Docker Compose
```bash
docker-compose up -d
```
This will start the application, the database, and the logging stack.

## API Documentation

The project includes an OpenAPI 3.0 specification. You can find it in the `openapi.yaml` file. You can use tools like Swagger UI or Redoc to visualize and interact with the API.

### Endpoints:
- `POST /register`: Register a new user.
- `POST /login`: Authenticate and receive tokens.
- `POST /refresh`: Refresh an expired access token.
- `POST /logout`: Revoke tokens and log out.
- `GET /protected`: Access a JWT-protected resource.
- `GET /ping`: Health check.

## Project Structure

- `cmd/app/`: Application entry point.
- `internal/api/`: HTTP handlers, routes, and DTOs.
- `internal/database/`: Database connection, repositories, and migrations.
- `internal/manager/`: JWT logic and password hashing.
- `internal/logger/`: Structured logging configuration and middleware.

## Note on Security
This project includes a `.env` file for convenience during development/learning with the provided `docker-compose.yml`. **Never** include sensitive credentials in version control for production environments. Always use secure secret management (e.g., HashiCorp Vault, AWS Secrets Manager).
