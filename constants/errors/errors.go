package errors

import "errors"

var (
	ErrInternalServer              = errors.New("internal server error")
	ErrUserExist                   = errors.New("user already exist")
	ErrUserNotFound                = errors.New("user not found")
	ErrInvalidInputBalance         = errors.New("balance must be positive")
	ErrInvalidEmailOrPasswordLogin = errors.New("invalid email or password")
	ErrCannotDeleteRootAdmin       = errors.New("cannot delete root admin")
)

var (
	ResUserCreated = map[string]string{"message": "User created"}
)
