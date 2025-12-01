package route

import (
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

	func InitRouter(app *gin.Engine, db *pgxpool.Pool, rd *redis.Client) {
		app.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(404, models.ResponseFailed{
				Success: false,
				Message: "Route Not Found, Try Again!",
			})
		})

	}
