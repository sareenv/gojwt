JWT implementation in Go using pgx and Postgres

Overview

This repository provides a minimal JWT-based authentication implementation written in Go. User credentials and session-related data are persisted in PostgreSQL using the pgx driver.

Features

- Sign-up and login endpoints
- JWT issuance and verification
- Token storage and revocation support via Postgres (pgx)

Prerequisites

- Go 1.20+ installed
- PostgreSQL accessible

Configuration

Set these environment variables before running:

- DATABASE_URL - Postgres connection string (eg. postgres://user:pass@host:5432/dbname)
- JWT_SECRET - secret key for signing tokens
- PORT - optional, default 8080

Running

go build ./... && ./your-binary

Or during development:

go run ./cmd/main.go

Usage

The server exposes endpoints for registration, login, and protected resources. On successful login a JWT is returned; include it in the Authorization header as "Bearer <token>" for protected requests.

Notes

The implementation uses github.com/jackc/pgx for database interactions and standard Go libraries for crypto and http handling.

License

MIT
