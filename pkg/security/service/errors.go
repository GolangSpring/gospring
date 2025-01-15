package service

import "errors"

var (
	UserNotFound = errors.New("UserFound")
	UserExists   = errors.New("UserExists")
)
