package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/config"
	"github.com/federus1105/koda-b4-final-backend/internal/middleware"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/federus1105/koda-b4-final-backend/internal/route"
	"github.com/federus1105/koda-b4-final-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())

	// --- LOAD .ENV IF DEVELOPMENT ---
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}
	router.Use(middleware.CORSMiddleware)

	// --- INIT DB ---
	db, err := config.InitDB()
	if err != nil {
		log.Println("Failed to connect to database\nCause: ", err.Error())
		return
	}
	defer db.Close()
	log.Println("DB Connected")

	// --- INIT RDB ---
	rdb, Rdb, err := config.InitRedis()
	if err != nil {
		log.Println("Failed to connect to redis\nCause: ", err.Error())
		return
	}
	defer rdb.Close()
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		fmt.Println("Failed Connected Redis : ", err.Error())
		return
	}
	
	// --- RATE LIMIT ---
	router.Use(middleware.RateLimiter(rdb, 100, 1*time.Minute))
	log.Println("REDIS Connected : ", Rdb)

	utils.StartCron(db)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, models.ResponseSucces{
			Success: true,
			Message: "Backend is running boy",
		})
	})

	route.InitRouter(router, db, rdb)
	router.Run(":8011")
}
