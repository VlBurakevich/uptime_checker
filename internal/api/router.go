package api

import (
	"time"
	"uptime-checker/internal/api/handlers"
	"uptime-checker/internal/api/middleware"
	services2 "uptime-checker/internal/api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, jwtSecret string, tokenTtl time.Duration) *gin.Engine {
	r := gin.Default()

	authService := &services2.AuthService{DB: db, JWTSecret: jwtSecret, TokenTTL: tokenTtl}
	authHandler := &handlers.AuthHandler{Service: authService}
	siteService := &services2.SiteService{DB: db}
	siteHandler := &handlers.SiteHandler{Service: siteService}

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret))
	{
		api.GET("/me", authHandler.GetMe)
		api.POST("/sites", siteHandler.Create)
		api.PUT("/sites/:id", siteHandler.Update)
		api.DELETE("/sites/:id", siteHandler.Delete)
		api.GET("/sites", siteHandler.List)
	}

	return r
}
