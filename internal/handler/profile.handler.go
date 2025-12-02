package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profile *repository.ProfileRepository
}

func NewProfileHandler(profile *repository.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{profile: profile}
}

func (p *ProfileHandler) Profile(ctx *gin.Context) {
	// --- GET USER IN CONTEXT ---
	userIDInterface, exists := ctx.Get(middleware.UserIDKey)
	if !exists {
		ctx.JSON(401, models.ResponseFailed{
			Success: false,
			Message: "Unauthorized: user not logged in",
		})
		return
	}

	var userID int
	switch v := userIDInterface.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		ctx.JSON(401, models.ResponseFailed{
			Success: false,
			Message: "Invalid user ID type in context",
		})
		return
	}

	// ---- LIMITS QUERY EXECUTION TIME ---
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	profile, err := p.profile.Profile(ctxTimeout, userID)
	if err != nil {
		fmt.Println("error :", err)
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "internal server error",
		})
		return
	}
	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "get Profile Succesfully",
		Results: profile,
	})
}
