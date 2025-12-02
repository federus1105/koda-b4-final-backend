package handler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/libs"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/federus1105/koda-b4-final-backend/internal/repository"
	"github.com/federus1105/koda-b4-final-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	auth *repository.AuthRepository
}

func NewAuthHandler(auth *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (a *AuthHandler) Register(ctx *gin.Context) {
	var req models.AuthRegister
	// --- VALIDATION ---
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var msgs []string
			for _, fe := range ve {
				msgs = append(msgs, utils.ErrorMessage(fe))
			}
			ctx.JSON(400, models.ResponseFailed{
				Success: false,
				Message: strings.Join(msgs, ", "),
			})
			return
		}

		ctx.JSON(400, models.ResponseFailed{
			Success: false,
			Message: "invalid JSON format",
		})
		return
	}

	// --- HASHING ---
	hashed, err := libs.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to hash password",
		})
		return
	}

	// ---- LIMITS QUERY EXECUTION TIME ---
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	newUser, err := a.auth.Register(ctxTimeout, hashed, req)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "Register Succesfully",
		Results: gin.H{
			"id":       newUser.Id,
			"fullname": newUser.Fullname,
			"email":    newUser.Email,
		},
	})
}

func (a *AuthHandler) Login(ctx *gin.Context) {
	var req models.AuthLogin
	// --- VALIDATION ---
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var msgs []string
			for _, fe := range ve {
				msgs = append(msgs, utils.ErrorMessage(fe))
			}
			ctx.JSON(400, models.ResponseFailed{
				Success: false,
				Message: strings.Join(msgs, ", "),
			})
			return
		}

		ctx.JSON(400, models.ResponseFailed{
			Success: false,
			Message: "invalid JSON format",
		})
		return
	}

	// ---- LIMITS QUERY EXECUTION TIME ---
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	user, err := a.auth.Login(ctxTimeout, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			ctx.JSON(401, models.ResponseFailed{
				Success: false,
				Message: "Nama atau Password salah",
			})
			return
		}
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "internal server errorr",
		})
		return
	}

	// --- VERIFICATION HASH PASSWORD
	ok, err := libs.VerifyPassword(req.Password, user.Password)
	if err != nil || !ok {
		ctx.JSON(401, models.ResponseFailed{
			Success: false,
			Message: "invalid email or password",
		})
		return
	}

	// --- GENERATE JWT TOKEN
	claims := libs.NewJWTClaims(user.Id, user.Role)
	jwtToken, err := claims.GenToken()
	if err != nil {
		fmt.Println("Internal Server Error.\nCause: ", err)
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "internal server errorrr",
		})
		return
	}

	// generate refresh token
	refreshToken := utils.GenerateShortCode(64)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := a.auth.SaveRefreshToken(ctxTimeout, user.Id, refreshToken, expiresAt); err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to save refresh token",
		})
		return
	}

	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "login successful",
		Results: gin.H{
			"token":         jwtToken,
			"refresh_token": refreshToken,
		},
	})

}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var input models.RefreshToken
	if err := ctx.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		ctx.JSON(400, models.ResponseFailed{
			Success: false,
			Message: "invalid request",
		})
		return
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userID, err := h.auth.ValidateRefreshToken(ctxTimeout, input.RefreshToken)
	if err != nil {
		ctx.JSON(401, models.ResponseFailed{
			Success: false,
			Message: "invalid refresh token",
		})
		return
	}

	// --- generate new access token ---
	claims := libs.NewJWTClaims(userID, "user")
	accessToken, err := claims.GenToken()
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to generate access token",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	var input models.RefreshToken

	if err := ctx.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		ctx.JSON(400, models.ResponseFailed{
			Success: false,
			Message: "refresh token is required",
		})
		return
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// revoke refresh token
	err := h.auth.RevokeRefreshToken(ctxTimeout, input.RefreshToken)
	if err != nil {
		ctx.JSON(500, models.ResponseFailed{
			Success: false,
			Message: "failed to revoke refresh token",
		})
		return
	}

	ctx.JSON(200, models.ResponseSucces{
		Success: true,
		Message: "logout successful",
	})
}
