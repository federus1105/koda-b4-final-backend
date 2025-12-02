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
	shortLinkRepo := repository.NewShortlinkRepository(db)
	ShortlinkHandler := handler.NewShortlinkHandler(shortLinkRepo, os.Getenv("BASE_URL"))

	shortlinkRoute.POST("", middleware.OptionalAuth(), ShortlinkHandler.CreateShortlink)
}
