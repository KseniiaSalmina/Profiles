package database

import "errors"

var ErrUserAlreadyExist = errors.New("user with this id is already exist")
var ErrNotUniqueUsername = errors.New("user with this username is already exist")
var ErrUserDoesNotExist = errors.New("user does not exist")
