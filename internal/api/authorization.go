package api

import (
	"net/http"

	"github.com/KseniiaSalmina/Profiles/internal/api/validation"
)

func (s *Server) authorization(r *http.Request) (bool, error) {
	authString, ok := r.Header["Authorization"]
	if !ok {
		return false, validation.ErrIncorrectAuth
	}

	username, password, err := validation.AuthString(authString[0])
	if err != nil {
		return false, err
	}

	user, err := s.storage.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	if err := validation.User(username, password+s.storage.ReturnSalt(), *user); err != nil {
		return false, err
	}

	return user.Admin, nil
}
