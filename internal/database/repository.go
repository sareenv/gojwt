package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository interface {
	UpsertToken(ctx context.Context, userID string, token string, expiresAt time.Time) error
	GetToken(ctx context.Context, userID string) (string, error)
	DeleteToken(ctx context.Context, userID string) error
}

type User struct {
	UserID       string    `json:"user_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type PostgresRefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRefreshTokenRepository(pool *pgxpool.Pool) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{pool: pool}
}

func (r *PostgresRefreshTokenRepository) UpsertToken(ctx context.Context, userID string, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE
		SET token = EXCLUDED.token, expires_at = EXCLUDED.expires_at, created_at = CURRENT_TIMESTAMP
	`
	_, err := r.pool.Exec(ctx, query, userID, token, expiresAt)
	return err
}

func (r *PostgresRefreshTokenRepository) GetToken(ctx context.Context, userID string) (string, error) {
	query := `SELECT token FROM refresh_tokens WHERE user_id = $1`
	var token string
	err := r.pool.QueryRow(ctx, query, userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *PostgresRefreshTokenRepository) DeleteToken(ctx context.Context, userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING user_id, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(&user.UserID, &user.CreatedAt, &user.UpdatedAt)
	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT user_id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	var user User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.UserID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT user_id, email, password_hash, created_at, updated_at FROM users WHERE user_id = $1`
	var user User
	err := r.pool.QueryRow(ctx, query, id).Scan(&user.UserID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
