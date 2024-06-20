package weeny

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

var PresetUnloggedErrors = []error{
	ErrVisited,
	ErrURLInvalid,
	ErrNoElemFound,
	ErrMaxDepth,
	ErrForbiddenDomain,
	ErrNotInteractable,
}

var (
	// errors from bot
	ErrNoElemFound     = errors.New("no element found")
	ErrNotInteractable = errors.New("elem not interactable")
)

var (
	// errors of control logics
	ErrForbiddenDomain = errors.New("forbidden domain")
	ErrMaxDepth        = errors.New("max depth limit reached")
	ErrVisited         = errors.New("url already visited")
	ErrURLInvalid      = errors.New("url invalid error")
)

func errVisited(msg string) error {
	return fmt.Errorf("%w: %s", ErrVisited, msg)
}

func errURLInvalid(msg string) error {
	return fmt.Errorf("%w: %s", ErrURLInvalid, msg)
}

// Echo warn if log happens
func (c *Crawler) Echo(err error) {
	if err != nil {
		c.logger.Warn("unhandled error", zap.Error(err))
	}
}

// Pie print if error not in ignored errors
func (c *Crawler) Pie(err error) {
	err = c.filterErrors(err)
	if err != nil {
		c.logger.Warn("unhandled error", zap.Error(err))
	}
}

func (c *Crawler) filterErrors(err error) error {
	errs := PresetUnloggedErrors
	if len(c.ignoredErrors) != 0 {
		errs = c.ignoredErrors
	}

	for _, e := range errs {
		if errors.Is(err, e) {
			return nil
		}
	}

	return err
}
