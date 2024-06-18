package weeny

import (
	"errors"
	"fmt"
)

var (
	ErrForbiddenDomain = errors.New("forbidden domain")

	ErrNoElemFound = errors.New("no element found")
	ErrMaxDepth    = errors.New("max depth limit reached")
	ErrVisited     = errors.New("visited error")
	ErrURLInvalid  = errors.New("url invalid error")
)

func errVisited(msg string) error {
	return fmt.Errorf("%w: %s", ErrVisited, msg)
}

func errURLInvalid(msg string) error {
	return fmt.Errorf("%w: %s", ErrURLInvalid, msg)
}
