package route

import (
	"os"

	"github.com/federus1105/koda-b4-final-backend/internal/handler"
	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitShortLinkRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	shortlinkRoute := router.Group("/api/v1/links")
	shortLinkRepo := repository.NewShortlinkRepository(db, rdb)
	ShortlinkHandler := handler.NewShortlinkHandler(shortLinkRepo, os.Getenv("BASE_URL"), rdb)

	shortlinkRoute.POST("", middleware.OptionalAuth(), ShortlinkHandler.CreateShortlink)
	shortlinkRoute.GET("", middleware.AuthMiddleware(), ShortlinkHandler.GetListLinks)
	shortlinkRoute.DELETE("/:shortcode", middleware.AuthMiddleware(), ShortlinkHandler.DeleteShortlink)
	shortlinkRoute.GET("/:shortcode", middleware.AuthMiddleware(), ShortlinkHandler.GetShortlinkDetail)

	redirectRouter := router.Group("/")
	redirectRouter.GET("/:shortcode", ShortlinkHandler.Redirect)
}
