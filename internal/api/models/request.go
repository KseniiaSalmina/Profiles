package models

type UserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}
