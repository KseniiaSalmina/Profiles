package validation

import (
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
)

func AuthString(auth string) (string, string, error) {
	base64Info, ok := strings.CutPrefix(auth, "Basic ")
	if !ok {
		return "", "", ErrIncorrectAuth
	}

	info, err := base64.RawStdEncoding.DecodeString(base64Info)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode authorization string: %w", err)
	}

	authData := strings.Split(string(info), ":")
	if len(authData) != 2 {
		return "", "", ErrIncorrectAuth
	}

	return authData[0], authData[1], nil
}

func User(username, password string, user models.User) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrIncorrectUserData
	}

	if username != user.Username {
		return ErrIncorrectUserData
	}

	return nil
}
