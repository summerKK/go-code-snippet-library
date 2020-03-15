package session

import "errors"

var (
	keyNotExistsInSession = errors.New("key not exists in session")
	sessionNotExists      = errors.New("session not exists")
)
