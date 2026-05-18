// Package manager provides tools for managing JWT-based authentication,
// supporting access tokens, refresh token rotation, and absolute session limits.
package manager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPair represents a pair of access and refresh tokens returned to the client.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// JWTClaims extends standard JWT registered claims with custom fields like UserID
// and AbsoluteExpiresAt for session control.
type JWTClaims struct {
	UserID string `json:"user_id"`
	// AbsoluteExpiresAt is the unix timestamp when the entire session must end.
	AbsoluteExpiresAt int64 `json:"abs_exp"`
	jwt.RegisteredClaims
}

// JWTManager handles the creation, validation, and refreshing of JWT tokens.
type JWTManager struct {
	Config *JWTConfig
}

// NewJWTManager creates a new instance of JWTManager with the provided configuration.
func NewJWTManager(config *JWTConfig) *JWTManager {
	return &JWTManager{
		Config: config,
	}
}

// GenerateToken creates a new short-lived access token for a specific user.
func (m *JWTManager) GenerateToken(userID string) (string, error) {
	// Implement token generation logic here using m.Config.SecretKey
	now := time.Now()
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.Config.AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.Config.SecretKey))
}

// RefreshToken validates an existing refresh token and generates a new pair (Access + Refresh).
// This implements Refresh Token Rotation, ensuring the session is extended but still
// bound by the original AbsoluteExpiresAt limit.
func (m *JWTManager) RefreshToken(tokens TokenPair, userID string) (pair *TokenPair, err error) {
	refreshToken := tokens.RefreshToken
	claims, err := m.GetClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	if claims.UserID != userID {
		return nil, fmt.Errorf("invalid user id")
	}

	// Check absolute expiration
	if time.Now().Unix() > claims.AbsoluteExpiresAt {
		return nil, fmt.Errorf("session expired (absolute limit reached)")
	}

	newAccessToken, err := m.GenerateToken(userID)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := m.GenerateRotatedRefreshToken(userID, claims.AbsoluteExpiresAt)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// GenerateRefreshToken creates the initial refresh token for a new session.
// It calculates the AbsoluteExpiresAt based on the MaxSessionDuration.
func (m *JWTManager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	absExp := now.Add(m.Config.MaxSessionDuration).Unix()
	return m.GenerateRotatedRefreshToken(userID, absExp)
}

// GenerateRotatedRefreshToken is a helper to create a refresh token with a specific absolute expiration.
func (m *JWTManager) GenerateRotatedRefreshToken(userID string, absoluteExpiresAt int64) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:            userID,
		AbsoluteExpiresAt: absoluteExpiresAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.Config.RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.Config.SecretKey))
}

// GetClaims parses and validates a token string and returns the custom JWTClaims.
func (m *JWTManager) GetClaims(token string) (*JWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(m.Config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}

// ValidateRefreshToken checks if a refresh token is valid and hasn't hit its absolute expiration.
func (m *JWTManager) ValidateRefreshToken(token string, userID string) (bool, error) {
	claims, err := m.GetClaims(token)
	if err != nil {
		return false, err
	}

	if claims.UserID != userID {
		return false, fmt.Errorf("invalid user id")
	}

	if time.Now().Unix() > claims.AbsoluteExpiresAt {
		return false, fmt.Errorf("session expired (absolute limit reached)")
	}

	return true, nil
}

// ValidateAccessToken checks if an access token is valid and returns its claims.
func (m *JWTManager) ValidateAccessToken(token string) (*JWTClaims, error) {
	// Implement token validation logic here using m.Config.SecretKey
	parsedToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(m.Config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
