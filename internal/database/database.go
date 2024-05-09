package database

import (
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"github.com/KseniiaSalmina/Profiles/internal/config"
)

type Database struct {
	mutex       sync.RWMutex
	users       []*User
	idIDX       map[string]*User
	usernameIDX map[string]*User
	deleteIDX   map[string]int
}

func NewDatabase(cfg config.Database, salt string) (*Database, error) {
	d := Database{
		users:       make([]*User, 0, 10),
		idIDX:       make(map[string]*User),
		usernameIDX: make(map[string]*User),
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword+salt), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to create first admin: %w", err)
	}

	firstUser := User{
		Email:    cfg.AdminEmail,
		Username: cfg.AdminUsername,
		PassHash: string(hashPass),
		Admin:    true,
	}

	if err := d.AddUser(firstUser); err != nil {
		return nil, fmt.Errorf("failed to first admin: %w", err)
	}

	return &d, nil
}

func (db *Database) AddUser(user User) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, ok := db.idIDX[user.ID]; ok {
		return ErrUserAlreadyExist
	}

	if _, ok := db.usernameIDX[user.Username]; ok {
		return ErrNotUniqueUsername
	}

	db.idIDX[user.ID] = &user
	db.usernameIDX[user.Username] = &user
	db.users = append(db.users, &user)
	db.deleteIDX[user.ID] = len(db.users) - 1

	return nil
}

func (db *Database) GetAllUsers(offset, limit int) []User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	result := make([]User, 0, limit)

	for i := offset + 1; offset <= offset+limit; offset++ {
		if i == len(db.users) {
			break
		}

		result = append(result, *db.users[i])
	}

	return result
}

func (db *Database) CountUsers() int {
	return len(db.users)
}

func (db *Database) GetUserByID(id string) (*User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	user, ok := db.idIDX[id]
	if !ok {
		return nil, ErrUserDoesNotExist
	}

	return user, nil
}

func (db *Database) GetUserByUsername(username string) (*User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	user, ok := db.usernameIDX[username]
	if !ok {
		return nil, ErrUserDoesNotExist
	}

	return user, nil
}

func (db *Database) ChangeUser(user User) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	oldUser, ok := db.idIDX[user.ID]
	if !ok {
		return ErrUserDoesNotExist
	}

	if user.Username != oldUser.Username {
		if _, ok = db.usernameIDX[user.Username]; ok {
			return ErrNotUniqueUsername
		}

		delete(db.usernameIDX, oldUser.Username)
		db.usernameIDX[user.Username] = oldUser
	}

	*oldUser = user

	return nil
}

func (db *Database) DeleteUser(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	index, ok := db.deleteIDX[id]
	if !ok {
		return ErrUserDoesNotExist
	}

	db.users = append(db.users[:index], db.users[index+1:]...)
	db.users[len(db.users)-1] = nil

	return nil
}
