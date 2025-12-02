package route

import (
	"github.com/federus1105/koda-b4-final-backend/internal/handler"
	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitProfileRouter(router *gin.Engine, db *pgxpool.Pool, rd *redis.Client) {
	profileRouter := router.Group("/api/v1")
	authRepository := repository.NewProfileRepository(db, rd)
	ProfileHandler := handler.NewProfileHandler(authRepository)

	profileRouter.GET("/profile", middleware.AuthMiddleware(), ProfileHandler.Profile)
}
