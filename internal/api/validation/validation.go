package validation

import (
	"encoding/base64"
	"fmt"
	"hash/crc64"
	"strings"

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
	passwordHash := crc64.Checksum([]byte(password), crc64.MakeTable(crc64.ISO))
	if username != user.Username || passwordHash != user.Password {
		return ErrIncorrectPassword
	}

	return nil
}

func Admin(username, password string, user models.User) error {
	if err := User(username, password, user); err != nil {
		return fmt.Errorf("failed to check admin status: %w", err)
	}

	if !user.Admin {
		return ErrIsNotAdmin
	}

	return nil
}
