package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/federus1105/koda-b4-final-backend/internal/libs"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const UserIDKey ctxKey = "id_users"

// --- FUNCTION SET DATA USER TO CONTEXT ---
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header missing or invalid",
			})
			return
		}

		// --- GET TOKEN ---
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ---- INSTANCE CLAIMS ---
		claims := &libs.Claims{}

		// --- VERIFY TOKEN ---
		if err := claims.VerifyToken(tokenString); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token: " + err.Error(),
			})
			return
		}

		// --- SAVE ID AND ROLE TO CONTEXT ---
		c.Set(UserIDKey, claims.ID)
		c.Set("role", claims.Role)

		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.ID)
		c.Request = c.Request.WithContext(ctx)

		// --- NEXT TO CONTROLLER ---
		c.Next()
	}
}

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
