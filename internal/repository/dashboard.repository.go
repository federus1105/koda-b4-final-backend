package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) CountLinks(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM shortlink").Scan(&count)
	return count, err
}

func (r *DashboardRepository) CountVisits(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM click").Scan(&count)
	return count, err
}

func (r *DashboardRepository) VisitsLast7Days(ctx context.Context) ([]int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DATE(created_at), COUNT(*) 
		FROM click
		WHERE created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chart := make([]int, 7)
	i := 0
	for rows.Next() {
		var date time.Time
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			return nil, err
		}
		chart[i] = count
		i++
	}
	return chart, nil
}
