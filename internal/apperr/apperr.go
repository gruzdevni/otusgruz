package apperr

import "errors"

var (
	ErrNotCorrectData   = errors.New("Not correct password or email")
	ErrNoSuchUser       = errors.New("No such email registered")
	ErrEmailAlreadyUsed = errors.New("Email is already registered. Please login")
	ErrNotEnoughFunds   = errors.New("Not enough funds to perform operation")
	ErrNoPermission     = errors.New("No permission to perform action")
)
