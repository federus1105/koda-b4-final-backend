package middleware

import (
	"context"
	"strings"

	"github.com/federus1105/koda-b4-final-backend/internal/libs"
	"github.com/gin-gonic/gin"
)

func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// --- IF NULL -> ANONYMOUS ---
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &libs.Claims{}
		if err := claims.VerifyToken(tokenString); err != nil {
			c.Next()
			return
		}

		c.Set(UserIDKey, claims.ID)
		c.Set("role", claims.Role)

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.ID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
