package models

type RegisterUserData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type LoginUserData struct {
	Email    string `json:"email"`
	Password string  `json:"password"`
}