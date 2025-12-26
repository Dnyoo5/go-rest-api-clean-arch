package models

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=3,max=20"`
}
