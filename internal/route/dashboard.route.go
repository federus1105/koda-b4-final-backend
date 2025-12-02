package route

import (
	"github.com/federus1105/koda-b4-final-backend/internal/handler"
	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitDashboardRouter(router *gin.Engine, db *pgxpool.Pool, rd *redis.Client) {
	dashboardRouter := router.Group("/api/v1/dashboard")
	dashboardRepository := repository.NewDashboardRepository(db)
	authHandler := handler.NewDashboardHandler(dashboardRepository, rd)

	dashboardRouter.GET("/stats", middleware.AuthMiddleware(), authHandler.Stats)
}
