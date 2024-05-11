package api

import (
	"errors"
	"net/http"

	"github.com/KseniiaSalmina/Profiles/internal/validation"
)

var ErrNoAuthString = errors.New("authorization required")

func (s *Server) authorization(r *http.Request) (bool, error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		return false, ErrNoAuthString
	}

	user, err := s.service.GetAuthData(username)
	if err != nil {
		return false, err
	}

	if err := validation.Auth(username, password+s.service.ReturnSalt(), *user); err != nil {
		return false, err
	}

	return user.Admin, nil
}
