package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/shortlink-backend/internal/controller"
	"github.com/shortlink-backend/internal/repository"
	"github.com/shortlink-backend/internal/service"
)

func AuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")

	authRepo := repository.NewAuthRepository(db)
	authCache := repository.NewAuthCacheRepository(rdb)
	// smtpMailer := pkg.NewSMTPMailerFromEnv()
	authService := service.NewAuthService(authRepo, authCache)
	authController := controller.NewAuthController(authService)

	// authRouter.POST("", authController.Login)
	authRouter.POST("/register", authController.Register)
	// authRouter.POST("/forgot-password", authController.ForgotPassword)
	// authRouter.POST("/reset-password", authController.ResetPassword)
	// authRouter.DELETE("/logout", authController.Logout)
}
