package repository

import (
	"context"
	"time"

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

func (sr *ShortlinkRepository) FindByCode(ctx context.Context, code string) (*models.Shortlink, error) {
	query := `
        SELECT id, account_id, original_url, is_active, expired_at
        FROM shortlink
        WHERE short_code = $1
    `

	row := sr.db.QueryRow(ctx, query, code)

	var s models.Shortlink
	err := row.Scan(&s.ID, &s.AccountID, &s.OriginalURL, &s.IsActive, &s.ExpiredAt)

	return &s, err
}

func (sr *ShortlinkRepository) InsertClick(ctx context.Context, shortlinkID int64, ip, userAgent, referer string) error {
	// --- START TRANSACTION ---
	tx, err := sr.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// --- INSERT TABLE CLICK ---
	_, err = tx.Exec(ctx, `
        INSERT INTO click(shortlink_id, ip, user_agent, referer)
        VALUES ($1, $2, $3, $4)
    `, shortlinkID, ip, userAgent, referer)
	if err != nil {
		return err
	}

	// --- UPDATE TOTAL CLICK ---
	_, err = tx.Exec(ctx, `
        UPDATE shortlink
        SET total_click = total_click + 1
        WHERE id = $1
    `, shortlinkID)
	if err != nil {
		return err
	}

	// --- COMMIT ---
	return tx.Commit(ctx)
}

func (r *ShortlinkRepository) DeactivateIfExpired(ctx context.Context, shortlink *models.Shortlink) (bool, error) {
	now := time.Now().UTC()
	expiredAt := shortlink.ExpiredAt.UTC()

	if now.After(expiredAt) && shortlink.IsActive {
		_, err := r.db.Exec(ctx, `
			UPDATE shortlink
			SET is_active = false
			WHERE id = $1
		`, shortlink.ID)
		if err != nil {
			return false, err
		}

		shortlink.IsActive = false
		return true, nil
	}

	return false, nil
}
