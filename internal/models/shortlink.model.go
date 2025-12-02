package models

import "time"

type CreateShortlinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	Expired     int    `json:"expired"`
}

type Shortlink struct {
	ID          int64      `json:"id"`
	AccountID   *int64     `json:"account_id"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	IsActive    bool       `json:"is_active"`
	ExpiredAt   *time.Time `json:"expired_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
