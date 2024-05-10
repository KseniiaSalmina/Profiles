package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/uptrace/bunrouter"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/api/validation"
)

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	var sc int
	defer s.logging(&sc, r)

	if _, err := s.authorization(r); err != nil {
		s.logger.WithError(err).Info("get all users handler, failed authorization")
		sc = http.StatusUnauthorized
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	limit, pageNo, err := s.getPageInfo(r)
	if err != nil {
		s.logger.WithError(err).Info("get all users handler, failed to get page info")
		sc = http.StatusBadRequest
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users := s.storage.GetAllUsers((pageNo-1)*limit, limit)

	sc = http.StatusOK
	_ = json.NewEncoder(w).Encode(users)
}

func (s *Server) postUser(w http.ResponseWriter, r *http.Request) {
	var sc int
	defer s.logging(&sc, r)

	isAdmin, err := s.authorization(r)
	if err != nil {
		s.logger.WithError(err).Info("post user handler, failed authorization")
		sc = http.StatusUnauthorized
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !isAdmin {
		s.logger.Info("post user handler, user is not admin")
		sc = http.StatusForbidden
		http.Error(w, validation.ErrIsNotAdmin.Error(), http.StatusForbidden)
		return
	}

	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.logger.WithError(err).Info("post user handler, failed unmarshall request body")
		sc = http.StatusInternalServerError
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := s.storage.AddUser(user)
	if err != nil {
		s.logger.WithError(err).Info("post user handler, failed to add user")
		sc = http.StatusBadRequest
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sc = http.StatusOK
	_ = json.NewEncoder(w).Encode(id)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	var sc int
	defer s.logging(&sc, r)

	if _, err := s.authorization(r); err != nil {
		s.logger.WithError(err).Info("get user handler, failed authorization")
		sc = http.StatusUnauthorized
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id, ok := bunrouter.ParamsFromContext(r.Context()).Get("id")
	if !ok {
		s.logger.Info("get user handler, failed to get id")
		sc = http.StatusBadRequest
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		s.logger.WithError(err).Info("get user handler, failed to parce uuid")
		sc = http.StatusBadRequest
		http.Error(w, "id should be in uuid format", http.StatusBadRequest)
		return
	}

	user, err := s.storage.GetUserByID(id)
	if err != nil {
		s.logger.WithError(err).Info("get user handler, failed to get user by id")
		sc = http.StatusBadRequest
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sc = http.StatusOK
	_ = json.NewEncoder(w).Encode(user)

}

func (s *Server) putUser(w http.ResponseWriter, r *http.Request) {
	var sc int
	defer s.logging(&sc, r)

	isAdmin, err := s.authorization(r)
	if err != nil {
		s.logger.WithError(err).Info("put user handler, failed authorization")
		sc = http.StatusUnauthorized
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !isAdmin {
		s.logger.Info("put user handler, user is not admin")
		sc = http.StatusForbidden
		http.Error(w, validation.ErrIsNotAdmin.Error(), http.StatusForbidden)
		return
	}

	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.logger.WithError(err).Info("put user handler, failed to unmarshall request body")
		sc = http.StatusInternalServerError
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = uuid.Parse(user.ID)
	if err != nil {
		s.logger.WithError(err).Info("put user handler, failed to parce uuid")
		sc = http.StatusBadRequest
		http.Error(w, "id should be in uuid format", http.StatusBadRequest)
		return
	}

	if err := s.storage.ChangeUser(user); err != nil {
		s.logger.WithError(err).Info("put user handler, failed to change user")
		sc = http.StatusBadRequest
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sc = http.StatusOK
	w.WriteHeader(http.StatusOK)

}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	var sc int
	defer s.logging(&sc, r)

	isAdmin, err := s.authorization(r)
	if err != nil {
		s.logger.WithError(err).Info("delete user handler, failed authorization")
		sc = http.StatusUnauthorized
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !isAdmin {
		s.logger.Info("delete user handler, user is not admin")
		sc = http.StatusForbidden
		http.Error(w, validation.ErrIsNotAdmin.Error(), http.StatusForbidden)
		return
	}

	id, ok := bunrouter.ParamsFromContext(r.Context()).Get("id")
	if !ok {
		s.logger.Info("delete user handler, failed to get id")
		sc = http.StatusBadRequest
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err = uuid.Parse(id)
	if err != nil {
		s.logger.WithError(err).Info("delete user handler, failed to parce uuid")
		sc = http.StatusBadRequest
		http.Error(w, "id should be in uuid format", http.StatusBadRequest)
		return
	}

	if err := s.storage.DeleteUser(id); err != nil {
		s.logger.WithError(err).Info("delete user handler, failed to delete user")
		sc = http.StatusBadRequest
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sc = http.StatusOK
	w.WriteHeader(http.StatusOK)
}
