package models

type DashboardStats struct {
	TotalLinks     int                      `json:"total_links"`
	TotalVisits    int                      `json:"total_visits"`
	AvgClickRate   float64                  `json:"avg_click_rate"`
	Last7DaysChart []int `json:"last_7_days_chart"`
}
