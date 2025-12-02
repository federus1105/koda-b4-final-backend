package models

type AuthRegister struct {
	Id              int    `json:"id"`
	Fullname        string `json:"fullname" binding:"required,max=20"`
	Email           string `json:"email"  binding:"required,email"`
	Password        string `json:"password" binding:"required,password_complex"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type AuthLogin struct {
	Id       int    `json:"id"`
	Email    string `json:"email"  binding:"required,email"`
	Password string `json:"password" binding:"required,password_complex"`
	Role     string `json:"role,omitempty"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
