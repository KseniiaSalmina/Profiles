package database

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/KseniiaSalmina/Profiles/internal/config"
)

type Database struct {
	mutex       sync.RWMutex
	users       []*User
	idIDX       map[string]*User
	usernameIDX map[string]*User
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
		ID:       uuid.NewString(),
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

	return nil
}

func (db *Database) GetAllUsers(offset, limit int) []User {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if offset > len(db.users)-1 {
		return []User{}
	}

	result := make([]User, 0, limit)

	from := offset
	to := offset + limit

	if to > len(db.users)-1 {
		for _, user := range db.users[from:] {
			result = append(result, *user)
		}
	} else {
		for _, user := range db.users[from:to] {
			result = append(result, *user)
		}
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

func (db *Database) ChangeUser(user UserUpdate) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	oldUser, ok := db.idIDX[user.ID]
	if !ok {
		return ErrUserDoesNotExist
	}

	if *user.Username != oldUser.Username {
		if _, ok = db.usernameIDX[*user.Username]; ok {
			return ErrNotUniqueUsername
		}

		delete(db.usernameIDX, oldUser.Username)
		db.usernameIDX[*user.Username] = oldUser
	}

	newUser := db.updateUser(*oldUser, user)

	*oldUser = newUser

	return nil
}

func (db *Database) updateUser(oldUser User, changes UserUpdate) User {
	newUser := User{
		ID: oldUser.ID,
	}

	if changes.Email != nil {
		newUser.Email = *changes.Email
	}

	if changes.Username != nil {
		newUser.Username = *changes.Username
	}

	if changes.PassHash != nil {
		newUser.PassHash = *changes.PassHash
	}

	if changes.Admin != nil {
		newUser.Admin = *changes.Admin
	}

	return newUser
}

func (db *Database) DeleteUser(id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	user, ok := db.idIDX[id]
	if !ok {
		return ErrUserDoesNotExist
	}

	delete(db.idIDX, user.ID)
	delete(db.usernameIDX, user.Username)

	for i, v := range db.users {
		if v.ID == user.ID {
			if i != len(db.users)-1 {
				db.users = append(db.users[:i], db.users[i+1:]...)
			}
			db.users = db.users[:len(db.users)-1]
			break
		}
	}

	return nil
}
