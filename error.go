package classify

import (
	"errors"
)

type ErrorClassify error

var (
	ErrGetKeyNotExists = errors.New("Get() the key is not exists")
	ErrKeysNotExists   = errors.New("Keys() the key is not exists")
)
