# JWT Authentication in Go

A minimal JWT authentication implementation in Go using PostgreSQL and `pgx`.

## Overview

This repository provides a simple authentication service with JWT-based sign-up, login, token verification, and token revocation support.

User credentials and session-related data are stored in PostgreSQL using the [`pgx`](https://github.com/jackc/pgx) driver.

## Features

- User sign-up and login
- JWT issuance and verification
- Token persistence and revocation
- PostgreSQL integration with `pgx`
- Minimal HTTP and crypto implementation using Go standard libraries

## Prerequisites

- Go 1.20+
- PostgreSQL

## Configuration

Set the following environment variables before running the application:

```sh
DATABASE_URL=postgres://user:pass@host:5432/dbname
JWT_SECRET=your-secret-key
PORT=8080
