package database

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrConnection = errors.New("connection error")

func Connect(ctx context.Context, cfg *DBConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnectionURL())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return nil, ErrConnection
	}
	slog.Info("connected to database")
	return pool, nil
}
