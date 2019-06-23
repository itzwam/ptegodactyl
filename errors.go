package ptegodactyl

import (
	"errors"
)

// ErrTooManyRequests is the error returned if you request too fast the API
var ErrTooManyRequests = errors.New("You're requesting too many much! Slow down")

// ErrUnauthorized is the error returned if you aren't Authorized
var ErrUnauthorized = errors.New("You didn't pass an auth header or it was missing the bearer")
