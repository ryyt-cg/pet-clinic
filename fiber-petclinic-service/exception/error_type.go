package exception

import "errors"

var (
	ErrorNotFound = errors.New("resource not found")
)
