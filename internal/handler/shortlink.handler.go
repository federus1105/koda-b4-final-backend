package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/federus1105/koda-b4-final-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ShortlinkHandler struct {
	repo    *repository.ShortlinkRepository
	rd      *redis.Client
	baseURL string
}

func NewShortlinkHandler(repo *repository.ShortlinkRepository, baseURL string, rd *redis.Client) *ShortlinkHandler {
	return &ShortlinkHandler{repo: repo, baseURL: baseURL, rd: rd}
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

func (h *ShortlinkHandler) Redirect(ctx *gin.Context) {
	code := ctx.Param("shortcode")

	// --- SEARCH LINK ---
	shortlink, err := h.repo.FindByCode(ctx, code)
	if err != nil {
		ctx.JSON(404, models.ResponseFailed{
			Success: false,
			Message: "Shortlink not found",
		})
		return
	}

	// --- CHECK EXPIRED ---
	expired, err := h.repo.DeactivateIfExpired(ctx, shortlink)
	if err != nil {
		log.Println("Failed to check expiration:", err)
	}
	if expired {
		ctx.JSON(410, models.ResponseFailed{
			Success: false,
			Message: "Shortlink is inactive",
		})
		return
	}

	err = h.repo.InsertClick(
		ctx,
		shortlink.ID,
		ctx.ClientIP(),
		ctx.Request.UserAgent(),
		ctx.Request.Referer(),
	)
	if err != nil {
		log.Println("Failed to record click:", err)
	}

	// --- REDIRECT ---
	ctx.Redirect(http.StatusFound, shortlink.OriginalURL)
}

func (h *ShortlinkHandler) GetListLinks(ctx *gin.Context) {
	// --- GET QUERY PARAMS ---
	pageStr := ctx.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit := 10
	offset := (page - 1) * limit

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

	// --- GET DATA WITH CACHE ---
	links, err := h.repo.GetListLinksByUser(ctxTimeout, h.rd, userID, limit, offset)
	if err != nil {
		log.Println("error getting user links:", err)
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "Get data successfully",
		Results: links,
	})
}

func (h *ShortlinkHandler) DeleteShortlink(ctx *gin.Context) {
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
	shortcode := ctx.Param("shortcode")
	if shortcode == "" {
		ctx.JSON(400, models.ResponseFailed{
			Success: false,
			Message: "Shortcode required",
		})
		return
	}

	// ---- LIMITS QUERY EXECUTION TIME ---
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.repo.DeleteShortlink(ctxTimeout, userID, shortcode); err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "Internal server error",
		})
		fmt.Println(err)
		return
	}

	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "Shortlink deleted successfully",
		Results: shortcode,
	})
}

func (h *ShortlinkHandler) GetShortlinkDetail(ctx *gin.Context) {
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

	shortcode := ctx.Param("shortcode")
	if shortcode == "" {
		ctx.JSON(400, models.ResponseFailed{
			Success: false,
			Message: "Shortcode required",
		})
		return
	}

	// ---- LIMITS QUERY EXECUTION TIME ---
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	link, err := h.repo.GetShortlinkDetail(ctxTimeout, userID, shortcode)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "Internal server error",
		})
		fmt.Println(err)
		return
	}

	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "get shortlink detail success",
		Results: link,
	})
}
