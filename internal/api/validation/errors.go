package validation

import "errors"

var ErrIncorrectAuth = errors.New("incorrect authorization string")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrIsNotAdmin = errors.New("user is not admin")
