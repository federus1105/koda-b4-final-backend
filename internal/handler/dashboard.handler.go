package handler

import (
	"fmt"
	"log"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/libs"
	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type DashboardHandler struct {
	dashboard *repository.DashboardRepository
	rd        *redis.Client
}

func NewDashboardHandler(dashboard *repository.DashboardRepository, rd *redis.Client) *DashboardHandler {
	return &DashboardHandler{dashboard: dashboard, rd: rd}
}

func (h *DashboardHandler) Stats(ctx *gin.Context) {

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
	key := fmt.Sprintf("analytics:%d:7d", userID)

	// --- GET CACHE ---
	stats, err := libs.GetFromCache[models.DashboardStats](ctx, h.rd, key)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to get cache",
		})
		return
	}

	if stats != nil {
		ctx.JSON(200, models.ResponseSucces{
			Success: true,
			Message: "Get from cache successfully",
			Results: stats,
		})
		return
	}

	// --- CACHE MISS ----
	totalLinks, err := h.dashboard.CountLinks(ctx)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to count links",
		})
		return
	}

	totalVisits, err := h.dashboard.CountVisits(ctx)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to count visits",
		})
		return
	}

	avgClickRate := 0.0
	if totalLinks > 0 {
		avgClickRate = float64(totalVisits) / float64(totalLinks)
	}

	chartData, err := h.dashboard.VisitsLast7Days(ctx)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to get last 7 days visits",
		})
		return
	}

	stats = &models.DashboardStats{
		TotalLinks:     totalLinks,
		TotalVisits:    totalVisits,
		AvgClickRate:   avgClickRate,
		Last7DaysChart: chartData,
	}

	// --- SAVE TO CACHE ---
	if err := libs.SetToCache(ctx, h.rd, key, stats, 10*time.Minute); err != nil {
		log.Println("failed to set cache:", err)
	}

	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "Get data successfullt",
		Results: stats,
	})

}
