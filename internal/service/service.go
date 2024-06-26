package service

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/KseniiaSalmina/Profiles/internal/database"
	"github.com/KseniiaSalmina/Profiles/internal/validation"
)

type Storage interface {
	GetUserByUsername(username string) (*database.User, error)
	GetAllUsers(offset, limit int) []database.User
	CountUsers() int
	AddUser(user database.User) error
	GetUserByID(id string) (*database.User, error)
	ChangeUser(user database.UserUpdate) error
	DeleteUser(id string) error
}

type Service struct {
	storage Storage
	salt    string
}

func NewService(cfg config.Service, storage Storage) (*Service, error) {
	service := Service{
		storage: storage,
		salt:    cfg.Salt,
	}

	firstUser := models.UserAdd{
		Email:    cfg.AdminEmail,
		Username: cfg.AdminUsername,
		Password: cfg.AdminPassword,
		Admin:    true,
	}

	if err := validation.UserAdd(firstUser); err != nil {
		return nil, fmt.Errorf("failed to add firs admin to db: %w", err)
	}

	if _, err := service.AddUser(firstUser); err != nil {
		return nil, fmt.Errorf("failed to add firs admin to db: %w", err)
	}

	return &service, nil
}

func (s *Service) ReturnSalt() string {
	return s.salt
}

func (s *Service) GetAuthData(username string) (*database.User, error) {
	user, err := s.storage.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get auth data: %w", err)
	}

	return user, nil
}

func (s *Service) GetAllUsers(limit, offset, pageNo int) *models.PageUsers {
	dbUsers := s.storage.GetAllUsers(offset, limit)

	users := make([]models.UserResponse, 0, len(dbUsers))
	for _, user := range dbUsers {
		u := models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			Admin:    user.Admin,
		}
		users = append(users, u)
	}

	usersAmount := s.storage.CountUsers()
	pagesAmount := usersAmount / limit
	if usersAmount%limit != 0 {
		pagesAmount++
	}

	if pageNo > pagesAmount {
		pageNo = pagesAmount
	}

	return &models.PageUsers{
		Users:       users,
		PageNo:      pageNo,
		Limit:       limit,
		PagesAmount: pagesAmount,
	}
}

func (s *Service) AddUser(user models.UserAdd) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password+s.salt), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	id := uuid.NewString()

	dbUser := database.User{
		ID:       id,
		Email:    user.Email,
		Username: user.Username,
		PassHash: string(hashPass),
		Admin:    user.Admin,
	}

	if err := s.storage.AddUser(dbUser); err != nil {
		return "", fmt.Errorf("failed to create new user: %w", err)
	}

	return id, nil
}

func (s *Service) GetUserByID(id string) (*models.UserResponse, error) {
	dbUser, err := s.storage.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	user := models.UserResponse{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Username: dbUser.Username,
		Admin:    dbUser.Admin,
	}

	return &user, nil
}

func (s *Service) ChangeUser(id string, user models.UserUpdate) error {
	dbUser := database.UserUpdate{
		ID:       id,
		Email:    user.Email,
		Username: user.Username,
		Admin:    user.Admin,
	}

	if user.Password != nil {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(*user.Password+s.salt), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to change user: %w", err)
		}
		newHash := string(hashPass)
		dbUser.PassHash = &newHash
	}

	if err := s.storage.ChangeUser(dbUser); err != nil {
		return fmt.Errorf("failed to change user: %w", err)
	}

	return nil
}

func (s *Service) DeleteUser(id string) error {
	if err := s.storage.DeleteUser(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
