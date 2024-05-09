package database

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/KseniiaSalmina/Profiles/internal/api/models"
	"github.com/KseniiaSalmina/Profiles/internal/config"
)

type Database struct {
	mutex           sync.RWMutex
	users           map[string]models.User
	uniqueUsernames map[string]string
	salt            string
}

func NewDatabase(cfg config.Database) (*Database, error) {
	d := Database{
		users:           make(map[string]models.User),
		uniqueUsernames: make(map[string]string),
		salt:            cfg.Salt,
	}

	firstUser := models.User{
		Email:    cfg.AdminEmail,
		Username: cfg.AdminUsername,
		Password: cfg.AdminPassword,
		Admin:    true,
	}

	if _, err := d.AddUser(firstUser); err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &d, nil
}

func (db *Database) ReturnSalt() string {
	return db.salt
}

func (db *Database) AddUser(user models.User) (string, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.users[user.ID]; ok {
		return "", ErrUserAlreadyExist
	}

	if _, ok := db.uniqueUsernames[user.Username]; ok {
		return "", ErrNotUniqueUsername
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password+db.salt), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	user.Password = string(hashPass)
	user.ID = uuid.NewString()

	db.uniqueUsernames[user.Username] = user.ID
	db.users[user.ID] = user

	return user.ID, nil
}

func (db *Database) GetAllUsers() []models.User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	result := make([]models.User, 0, len(db.users))

	for _, user := range db.users {
		result = append(result, user)
	}

	return result
}

func (db *Database) GetUserByID(id string) (*models.User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	user, ok := db.users[id]
	if !ok {
		return nil, ErrUserDoesNotExist
	}

	return &user, nil
}

func (db *Database) GetUserByUsername(username string) (*models.User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	userID, ok := db.uniqueUsernames[username]
	if !ok {
		return nil, ErrUserDoesNotExist
	}

	user, ok := db.users[userID]
	if !ok {
		return nil, ErrUserDoesNotExist
	}

	return &user, nil
}

func (db *Database) ChangeUser(user models.User) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	oldUser, ok := db.users[user.ID]
	if !ok {
		return ErrUserDoesNotExist
	}

	if user.Username != oldUser.Username {
		if _, ok = db.uniqueUsernames[user.Username]; ok {
			return ErrNotUniqueUsername
		}

		delete(db.uniqueUsernames, oldUser.Username)
		db.uniqueUsernames[user.Username] = user.ID
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password+db.salt), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.Password = string(hashPass)

	db.users[user.ID] = user

	return nil
}

func (db *Database) DeleteUser(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user, ok := db.users[id]
	if !ok {
		return ErrUserDoesNotExist
	}

	delete(db.users, id)
	delete(db.uniqueUsernames, user.Username)

	return nil
}
