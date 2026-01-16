package domain

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNotRegistered = errors.New("user not registered")
	ErrInvalidCredentials = errors.New("wrong username or password")
	ErrDb = errors.New("some database related error")
	ErrInternalServer = errors.New("server error")
	ErrNotFound = errors.New("Requested resource not found")
)


