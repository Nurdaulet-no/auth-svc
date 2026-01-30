package domain

import  "errors"

var(
	ErrInvalidCreadiantals = errors.New("Invalid credentials")
	ErrEmailTaken = errors.New("Email is already taken")
	ErrUserNotFound = errors.New("User not found")
)

