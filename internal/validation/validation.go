package validation

import (
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/database"
)

func Auth(username, password string, user database.User) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password)); err != nil {
		return ErrIncorrectAuthData
	}

	if username != user.Username {
		return ErrIncorrectAuthData
	}

	return nil
}

func UserAdd(user models.UserAdd) error {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return ErrIncorrectUserData
	}

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	return nil
}

func UserUpdate(user models.UserUpdate) error {
	if user.Email == nil && user.Username == nil && user.Password == nil && user.Admin == nil {
		return ErrNoChanges
	}

	if user.Email != nil {
		if *user.Email == "" {
			return ErrIncorrectUserData
		}
		if _, err := mail.ParseAddress(*user.Email); err != nil {
			return fmt.Errorf("invalid email address: %w", err)
		}
	}

	if user.Username != nil && *user.Username == "" {
		return ErrIncorrectUserData
	}

	if user.Password != nil && *user.Password == "" {
		return ErrIncorrectUserData
	}

	return nil
}
