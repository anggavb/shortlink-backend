package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/shortlink-backend/internal/controller"
	"github.com/shortlink-backend/internal/middleware"
	"github.com/shortlink-backend/internal/repository"
	"github.com/shortlink-backend/internal/service"
)

func LinkRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	linkRouter := router.Group("/links")

	linkRepo := repository.NewLinkRepository(db)
	linkCache := repository.NewLinkCacheRepository(rdb)
	linkService := service.NewLinkService(linkRepo, linkCache)
	linkController := controller.NewLinkController(linkService)

	linkRouter.Use(middleware.VerifyToken())
	linkRouter.POST("", linkController.CreateLink)
	linkRouter.GET("", linkController.ListLinks)
	linkRouter.DELETE("/:id", linkController.DeleteLink)
}
