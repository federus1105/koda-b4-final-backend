package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/libs"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rd *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		identity := strings.ReplaceAll(c.ClientIP(), ":", "")
		endpoint := c.FullPath()
		key := "ratelimit:" + identity + ":" + endpoint

		counter, err := libs.GetFromCache[int](ctx, rd, key)
		if err != nil {
			c.JSON(500, models.ResponseFailed{
				Success: false,
				Message: "Redis error",
			})
			c.Abort()
			return
		}

		if counter == nil {
			_ = libs.SetToCache(ctx, rd, key, 1, window)
		} else if *counter >= limit {
			c.Header("Retry-After", window.String())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Too many requests",
			})
			c.Abort()
			return
		} else {
			*counter++
			_ = libs.SetToCache(ctx, rd, key, *counter, window)
		}

		c.Next()
	}
}
