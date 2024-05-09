package formatter

import (
	"fmt"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/config"
	"github.com/KseniiaSalmina/Profiles/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	GetUserByUsername(username string) (*database.User, error)
	GetAllUsers(offset, limit int) []database.User
	CountUsers() int
	AddUser(user database.User) error
	GetUserByID(id string) (*database.User, error)
	ChangeUser(user database.User) error
	DeleteUser(id string) error
}

type Formatter struct {
	storage Storage
	salt    string
}

func NewFormatter(cfg config.Formatter, storage Storage) *Formatter {
	return &Formatter{
		storage: storage,
		salt:    cfg.Salt,
	}
}

func (f *Formatter) ReturnSalt() string {
	return f.salt
}

func (f *Formatter) GetPasswordHashByUsername(username string) (string, error) {
	user, err := f.storage.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("failed to get auth data: %w", err)
	}

	return user.PassHash, nil
}

func (f *Formatter) GetAllUsers(offset, limit int) *models.PageUsers {
	dbUsers := f.storage.GetAllUsers(offset, limit)

	users := make([]models.UserResponse, 0, len(dbUsers))
	for _, user := range dbUsers {
		u := models.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		}
		users = append(users, u)
	}

	usersAmount := f.storage.CountUsers()
	pagesAmount := usersAmount / limit
	if usersAmount%limit != 0 {
		pagesAmount++
	}

	return &models.PageUsers{
		Users:       users,
		PageNo:      offset / limit,
		Limit:       limit,
		PagesAmount: pagesAmount,
	}
}

func (f *Formatter) AddUser(user models.UserRequest) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password+f.salt), bcrypt.DefaultCost)
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

	if err := f.storage.AddUser(dbUser); err != nil {
		return "", fmt.Errorf("failed to create new user: %w", err)
	}

	return id, nil
}

func (f *Formatter) GetUserByID(id string) (*models.UserResponse, error) {
	dbUsers, err := f.storage.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	user := models.UserResponse{
		ID:       dbUsers.ID,
		Email:    dbUsers.Email,
		Username: dbUsers.Username,
	}

	return &user, nil
}

func (f *Formatter) ChangeUser(user models.UserRequest) error {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password+f.salt), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to change user: %w", err)
	}

	dbUser := database.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		PassHash: string(hashPass),
		Admin:    user.Admin,
	}

	if err := f.storage.ChangeUser(dbUser); err != nil {
		return fmt.Errorf("failed to change user: %w", err)
	}

	return nil
}

func (f *Formatter) DeleteUser(id string) error {
	if err := f.storage.DeleteUser(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
