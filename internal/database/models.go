package database

type User struct {
	ID       string
	Email    string
	Username string
	PassHash string
	Admin    bool
}

type UserUpdate struct {
	ID       string
	Email    *string
	Username *string
	PassHash *string
	Admin    *bool
}
