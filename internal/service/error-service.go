package service

import "errors"

var (
	ErrInvalidURL    = errors.New("Invalid URL")
	ErrInvalidScheme = errors.New("Invalid Scheme")
	ErrInvalidHost   = errors.New("Invalid Host")
)
