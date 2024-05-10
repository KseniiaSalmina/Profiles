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
	if _, err := s.authorization(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	limit, pageNo, err := s.getPageInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users := s.storage.GetAllUsers((pageNo-1)*limit, limit)

	_ = json.NewEncoder(w).Encode(users)
}

func (s *Server) postUser(w http.ResponseWriter, r *http.Request) {
	isAdmin, err := s.authorization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !isAdmin {
		http.Error(w, validation.ErrIsNotAdmin.Error(), http.StatusForbidden)
		return
	}

	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := s.storage.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(id)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	if _, err := s.authorization(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id, ok := bunrouter.ParamsFromContext(r.Context()).Get("id")
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "id should be in uuid format", http.StatusBadRequest)
		return
	}

	user, err := s.storage.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(user)

}

func (s *Server) putUser(w http.ResponseWriter, r *http.Request) {
	isAdmin, err := s.authorization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !isAdmin {
		http.Error(w, validation.ErrIsNotAdmin.Error(), http.StatusForbidden)
		return
	}

	var user models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = uuid.Parse(user.ID)
	if err != nil {
		http.Error(w, "id should be in uuid format", http.StatusBadRequest)
		return
	}

	if err := s.storage.ChangeUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	isAdmin, err := s.authorization(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !isAdmin {
		http.Error(w, validation.ErrIsNotAdmin.Error(), http.StatusForbidden)
		return
	}

	id, ok := bunrouter.ParamsFromContext(r.Context()).Get("id")
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, "id should be in uuid format", http.StatusBadRequest)
		return
	}

	if err := s.storage.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
