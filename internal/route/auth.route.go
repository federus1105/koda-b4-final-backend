package route

import (
	"github.com/federus1105/koda-b4-final-backend/internal/handler"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/api/v1/auth")
	authRepository := repository.NewAuthRepository(db)
	authHandler := handler.NewAuthHandler(authRepository)

	authRouter.POST("/register", authHandler.Register)
	authRouter.POST("/login", authHandler.Login)
}
