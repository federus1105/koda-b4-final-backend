package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/federus1105/koda-b4-final-backend/internal/libs"
	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ProfileRepository struct {
	db *pgxpool.Pool
	rd *redis.Client
}

func NewProfileRepository(db *pgxpool.Pool, rd *redis.Client) *ProfileRepository {
	return &ProfileRepository{db: db, rd: rd}
}

func (p *ProfileRepository) Profile(ctx context.Context, userId int) (*models.Profiles, error) {
	cacheKey := fmt.Sprintf("user:%d:profile", userId)

	// --- GET CACHE ---
	cached, err := libs.GetFromCache[models.Profiles](ctx, p.rd, cacheKey)
	if err != nil {
		log.Println("Redis GET error:", err)
	}
	if cached != nil {
		return cached, nil
	}

	var profile models.Profiles
	sql := `SELECT a.id, a.fullname, 
	a.photos, u.email FROM account a
	JOIN users u ON u.id = a.id_users
	WHERE u.id = $1`

	err = p.db.QueryRow(ctx, sql, userId).Scan(
		&profile.Id,
		&profile.Fullname,
		&profile.Photos,
		&profile.Email,
	)

	if err != nil {
		return &profile, err
	}
	// --- SAVE CACHE ----
	ttl := 24 * time.Hour
	if err := libs.SetToCache(ctx, p.rd, cacheKey, profile, ttl); err != nil {
		log.Println("Redis SET error:", err)
	}

	return &profile, nil
}
