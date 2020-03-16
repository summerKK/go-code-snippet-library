package session

import "errors"

var (
	errKeyNotExistsInSession = errors.New("key not exists in session")
	errSessionNotExists      = errors.New("session not exists")
)
