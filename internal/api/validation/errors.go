package validation

import "errors"

var ErrIncorrectAuth = errors.New("incorrect authorization string")
var ErrIncorrectUserData = errors.New("user with this username or password is not exist")
var ErrIsNotAdmin = errors.New("user is not admin")
