package validation

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/KseniiaSalmina/Profiles/internal/database"
)

func User(username, password string, user database.User) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		return ErrIncorrectUserData
	}

	if username != user.Username {
		return ErrIncorrectUserData
	}

	return nil
}
