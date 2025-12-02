package handler

import (
	"fmt"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/federus1105/koda-b4-final-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type ShortlinkHandler struct {
	repo    *repository.ShortlinkRepository
	baseURL string
}

func NewShortlinkHandler(repo *repository.ShortlinkRepository, baseURL string) *ShortlinkHandler {
	return &ShortlinkHandler{repo: repo, baseURL: baseURL}
}

func (h *ShortlinkHandler) CreateShortlink(c *gin.Context) {

	var req models.CreateShortlinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, models.ResponseFailed{
			Success: false,
			Message: "Invalid URL",
		})
		return
	}

	// OPTIONAL USER ID
	var accountID *int64 = nil
	if idVal, ok := c.Get(middleware.UserIDKey); ok {
		uid := int64(idVal.(int))
		accountID = &uid
	}

	// --- EXPIRE ---
	var expiredAt *time.Time
	if req.Expired > 0 {
		t := time.Now().Add(time.Duration(req.Expired) * 24 * time.Hour)
		expiredAt = &t
	} else {
		t := time.Now().Add(30 * 24 * time.Hour)
		expiredAt = &t
	}

	shortCode := utils.GenerateShortCode(6)
	shortlink := &models.Shortlink{
		AccountID:   accountID,
		ShortCode:   shortCode,
		OriginalURL: req.OriginalURL,
		IsActive:    true,
		ExpiredAt:   expiredAt,
	}
	if err := h.repo.CreateShortlink(c, shortlink); err != nil {
		c.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "Failed to create shortlink",
		})
		fmt.Println(err)
		return
	}

	fullShortURL := fmt.Sprintf("%s/%s", h.baseURL, shortCode)
	c.JSON(201, models.ResponseSucces{
		Success: true,
		Message: "create shortlink Succesfully",
		Results: gin.H{
			"id":           shortlink.ID,
			"original_url": shortlink.OriginalURL,
			"short_url":    fullShortURL,
			"expired_at":   shortlink.ExpiredAt,
			"created_at":   shortlink.CreatedAt,
		},
	})
}
