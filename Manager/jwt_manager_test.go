package manager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager(t *testing.T) {
	config := &JWTConfig{
		SecretKey:            "test_secret_key",
		AccessTokenDuration:  time.Minute * 1,
		RefreshTokenDuration: time.Hour * 1,
		MaxSessionDuration:   time.Hour * 24,
	}
	jwtManager := NewJWTManager(config)

	userID := "test-user-123"

	t.Run("GenerateAndValidateAccessToken", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(userID)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := jwtManager.ValidateAccessToken(token)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})

	t.Run("GenerateAndValidateRefreshToken", func(t *testing.T) {
		token, err := jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		valid, err := jwtManager.ValidateRefreshToken(token, userID)
		require.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("RefreshToken", func(t *testing.T) {
		refreshToken, err := jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)

		pair := TokenPair{
			RefreshToken: refreshToken,
		}

		newPair, err := jwtManager.RefreshToken(pair, userID)
		require.NoError(t, err)
		assert.NotEmpty(t, newPair.AccessToken)
		assert.NotEmpty(t, newPair.RefreshToken)

		// Validate new access token
		claims, err := jwtManager.ValidateAccessToken(newPair.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)

		// Validate new refresh token
		valid, err := jwtManager.ValidateRefreshToken(newPair.RefreshToken, userID)
		require.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("ExpiredAccessToken", func(t *testing.T) {
		shortConfig := &JWTConfig{
			SecretKey:            "test_secret_key",
			AccessTokenDuration:  -time.Second, // Already expired
			RefreshTokenDuration: time.Hour,
		}
		shortManager := NewJWTManager(shortConfig)

		token, err := shortManager.GenerateToken(userID)
		require.NoError(t, err)

		claims, err := shortManager.ValidateAccessToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("InvalidSecret", func(t *testing.T) {
		token, err := jwtManager.GenerateToken(userID)
		require.NoError(t, err)

		wrongConfig := &JWTConfig{
			SecretKey:            "wrong_secret",
			AccessTokenDuration:  time.Minute,
			RefreshTokenDuration: time.Hour,
		}
		wrongManager := NewJWTManager(wrongConfig)

		claims, err := wrongManager.ValidateAccessToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("InvalidUserIDInRefresh", func(t *testing.T) {
		refreshToken, err := jwtManager.GenerateRefreshToken(userID)
		require.NoError(t, err)

		valid, err := jwtManager.ValidateRefreshToken(refreshToken, "wrong-user")
		assert.Error(t, err)
		assert.False(t, valid)
		assert.Contains(t, err.Error(), "invalid user id")
	})

	t.Run("AbsoluteExpiration", func(t *testing.T) {
		absConfig := &JWTConfig{
			SecretKey:            "test_secret_key",
			AccessTokenDuration:  time.Minute,
			RefreshTokenDuration: time.Hour,
			MaxSessionDuration:   -time.Second, // Already expired absolute session
		}
		absManager := NewJWTManager(absConfig)

		refreshToken, err := absManager.GenerateRefreshToken(userID)
		require.NoError(t, err)

		// ValidateRefreshToken should fail
		valid, err := absManager.ValidateRefreshToken(refreshToken, userID)
		assert.Error(t, err)
		assert.False(t, valid)
		assert.Contains(t, err.Error(), "absolute limit reached")

		// RefreshToken should fail
		pair := TokenPair{RefreshToken: refreshToken}
		newPair, err := absManager.RefreshToken(pair, userID)
		assert.Error(t, err)
		assert.Nil(t, newPair)
		assert.Contains(t, err.Error(), "absolute limit reached")
	})
}
