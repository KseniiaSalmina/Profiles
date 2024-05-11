package validation

import "errors"

var ErrIncorrectUserData = errors.New("user with this username or password is not exist")
var ErrIsNotAdmin = errors.New("user is not admin")
