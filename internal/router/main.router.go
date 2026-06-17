package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	_ "github.com/shortlink-backend/docs"
	"github.com/shortlink-backend/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.Use(middleware.CORSMiddleware)
	router.Static("/img", "./public/img")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouter(router, db, rdb)
}
