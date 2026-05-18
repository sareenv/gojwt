package manager

import (
	"time"

	"github.com/caarlos0/env/v6"
)

// JWTConfig holds the configuration settings for JWT generation and validation.
// It uses environment variables for configuration with sensible defaults.
type JWTConfig struct {
	// SecretKey is the signing key used for HS256.
	SecretKey string `env:"JWT_SECRET_KEY" envDefault:"mysecretkey"`
	// AccessTokenDuration defines how long a short-lived access token is valid (e.g., 15m).
	AccessTokenDuration time.Duration `env:"ACCESS_TOKEN_DURATION" envDefault:"15m"`
	// RefreshTokenDuration defines the sliding window for a single refresh token (e.g., 7 days).
	RefreshTokenDuration time.Duration `env:"REFRESH_TOKEN_DURATION" envDefault:"168h"`
	// MaxSessionDuration defines the absolute maximum age of a session (e.g., 30 days).
	// After this duration, the user must re-authenticate regardless of token rotation.
	MaxSessionDuration time.Duration `env:"MAX_SESSION_DURATION" envDefault:"720h"`
}

// LoadJWTConfig parses environment variables into the JWTConfig struct.
func LoadJWTConfig() (*JWTConfig, error) {
	var cfg JWTConfig
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
