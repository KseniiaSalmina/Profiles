package database

import (
	"sync"
)

type Database struct {
	mutex       sync.RWMutex
	users       []*User
	idIDX       map[string]*User
	usernameIDX map[string]*User
}

func NewDatabase() *Database {
	return &Database{
		users:       make([]*User, 0),
		idIDX:       make(map[string]*User),
		usernameIDX: make(map[string]*User),
	}
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

	result := make([]User, 0, limit)

	from := offset
	to := offset + limit

	if offset > len(db.users)-1 {
		from = len(db.users) - limit
		to = len(db.users)
	}

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
	db.mutex.RLock()
	defer db.mutex.RUnlock()

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

	if user.Username != nil && *user.Username != oldUser.Username {
		if _, ok = db.usernameIDX[*user.Username]; ok {
			return ErrNotUniqueUsername
		}

		delete(db.usernameIDX, oldUser.Username)
		db.usernameIDX[*user.Username] = oldUser
	}

	db.updateUser(oldUser, user)

	return nil
}

func (db *Database) updateUser(user *User, changes UserUpdate) {
	if changes.Email != nil {
		user.Email = *changes.Email
	}

	if changes.Username != nil {
		user.Username = *changes.Username
	}

	if changes.PassHash != nil {
		user.PassHash = *changes.PassHash
	}

	if changes.Admin != nil {
		user.Admin = *changes.Admin
	}
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
