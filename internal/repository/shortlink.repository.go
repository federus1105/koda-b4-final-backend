package repository

import (
	"context"

	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortlinkRepository struct {
	db *pgxpool.Pool
}

func NewShortlinkRepository(db *pgxpool.Pool) *ShortlinkRepository {
	return &ShortlinkRepository{db: db}
}

func (r *ShortlinkRepository) CreateShortlink(ctx context.Context, s *models.Shortlink) error {
	query := `
        INSERT INTO shortlink (account_id, short_code, expired_at, original_url)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at;
    `

	return r.db.QueryRow(
		ctx,
		query,
		s.AccountID,
		s.ShortCode,
		s.ExpiredAt,
		s.OriginalURL,
	).Scan(&s.ID, &s.CreatedAt)
}