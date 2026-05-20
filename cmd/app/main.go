package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sareenv/gojwt/internal/api"
	"github.com/sareenv/gojwt/internal/database"
	"github.com/sareenv/gojwt/internal/logger"
	"github.com/sareenv/gojwt/internal/manager"
)

const (
	Port = ":8080"
)

func StartServer(jwtManager *manager.JWTManager, userRepo database.UserRepository, passwordHasher manager.PasswordHasher) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logger.Middleware())

	jwtHandler := api.NewJWTHandler(jwtManager, userRepo, passwordHasher)
	api.RegisterRoutes(router, jwtHandler)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := router.Run(Port)
	if err != nil {
		slog.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}

func main() {
	// Initialize logger
	logFile, err := logger.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Load environment variables from .env file.
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, relying on system environment variables")
	}

	// Load database configuration.
	dbCfg, err := database.LoadDBConfig()
	if err != nil {
		slog.Error("failed to load database config", "error", err)
		os.Exit(1)
	}
	slog.Info("Database configuration loaded", "config", dbCfg)

	pool, err := database.Connect(context.Background(), dbCfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Load JWT configuration.
	jwtCfg, err := manager.LoadJWTConfig()
	if err != nil {
		slog.Error("failed to load jwt config", "error", err)
		os.Exit(1)
	}

	repo := database.NewPostgresRefreshTokenRepository(pool)
	userRepo := database.NewPostgresUserRepository(pool)
	passwordHasher := manager.NewPasswordHasher()
	jwtManager := manager.NewJWTManager(jwtCfg, repo)

	StartServer(jwtManager, userRepo, passwordHasher)
}
