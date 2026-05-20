package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, handler *JWTHandler) {
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.POST("/refresh", handler.Refresh)
	router.POST("/logout", handler.Logout)
	router.GET("/protected", handler.Protected)
}
