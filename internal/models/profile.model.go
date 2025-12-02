package models

type Profiles struct {
	Id       int     `json:"id"`
	Fullname string  `json:"fullname"`
	Photos   *string `json:"photos"`
	Email    string  `json:"email"`
}
