package models

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Admin    bool   `json:"admin"`
}

type PageUsers struct {
	Users       []UserResponse `json:"users"`
	PageNo      int            `json:"page_number"`
	Limit       int            `json:"limit"`
	PagesAmount int            `json:"pages_amount"`
}
