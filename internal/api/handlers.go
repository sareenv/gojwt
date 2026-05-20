package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sareenv/gojwt/internal/database"
	"github.com/sareenv/gojwt/internal/manager"
)

type JWTHandler struct {
	jwtManager     *manager.JWTManager
	userRepo       database.UserRepository
	passwordHasher manager.PasswordHasher
}

func NewJWTHandler(jwtManager *manager.JWTManager, userRepo database.UserRepository, passwordHasher manager.PasswordHasher) *JWTHandler {
	return &JWTHandler{
		jwtManager:     jwtManager,
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (h *JWTHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid email or password"})
		return
	}

	valid, err := h.passwordHasher.Verify(req.Password, user.PasswordHash)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid email or password"})
		return
	}

	accessToken, err := h.jwtManager.GenerateToken(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate access token"})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(c.Request.Context(), user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *JWTHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	tokens := manager.TokenPair{
		RefreshToken: req.RefreshToken,
	}

	pair, err := h.jwtManager.RefreshToken(c.Request.Context(), tokens, req.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	})
}

func (h *JWTHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.jwtManager.Logout(c.Request.Context(), req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "successfully logged out"})
}

func (h *JWTHandler) Protected(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "authorization header missing"})
		return
	}

	// Expecting "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid authorization header format"})
		return
	}

	tokenString := authHeader[len(bearerPrefix):]
	claims, err := h.jwtManager.ValidateAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, ProtectedResponse{
		Message: "access granted",
		UserID:  claims.UserID,
	})
}

func (h *JWTHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPassword, err := h.passwordHasher.Hash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
		return
	}

	user := &database.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.userRepo.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}
