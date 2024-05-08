package models

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password uint64 `json:"password"`
	Admin    bool   `json:"admin"`
}
