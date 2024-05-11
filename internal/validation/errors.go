package validation

import "errors"

var ErrIncorrectAuthData = errors.New("user with this username or password is not exist")
var ErrIsNotAdmin = errors.New("user is not admin")
var ErrNoChanges = errors.New("no changes submitted")
var ErrIncorrectUserData = errors.New("user should have username, password and email")
