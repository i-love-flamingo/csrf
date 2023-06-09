package interfaces

import (
	"errors"
)

var (
	ErrInvalidToken = errors.New("csrf_token is not valid")
)
