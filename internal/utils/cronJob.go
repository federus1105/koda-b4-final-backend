package utils

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartCron(db *pgxpool.Pool) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			db.Exec(context.Background(),
				`UPDATE shortlink 
                 SET is_active = FALSE
                 WHERE expired_at IS NOT NULL 
                   AND expired_at < NOW()
                   AND is_active = TRUE`)
		}
	}()
}
