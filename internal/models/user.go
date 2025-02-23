package models

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=25"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}
