package weeny

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"time"

	"github.com/coghost/wee"
	"go.uber.org/zap"
)

func (c *Crawler) waitPageStable() {
	defer c.LogTimeSpent(time.Now())

	diff := c.stableDiff
	if err := c.Bot.Page().Timeout(wee.MediumToSec*time.Second).WaitDOMStable(time.Second*wee.NapToSec, diff); err != nil {
		c.logger.Warn("cannot wait dom stable", zap.Error(err))
	}
}

func (c *Crawler) ClosePage(err error) error {
	if errors.Is(err, ErrNoElemFound) {
		err = nil
	}

	if err != nil {
		return err
	}

	err = c.Bot.ResetToOriginalPage()
	if err != nil {
		return err
	}

	c.waitPageStable()

	return nil
}

func (c *Crawler) GoBack(err error) error {
	if errors.Is(err, ErrNoElemFound) {
		err = nil
	}

	if err != nil {
		return fmt.Errorf("visit elem failed: %w", err)
	}

	func() {
		defer c.LogTimeSpent(time.Now())
		c.Bot.Page().MustNavigateBack()
	}()

	c.waitPageStable()
	return nil
}

func (c *Crawler) LogTimeSpent(start time.Time) {
	const skip = 3

	if c.trackTime {
		TimeTrack(start, skip)
	}
}

func TimeTrack(start time.Time, skip int) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(skip) //nolint
	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	log.Printf("%s took %s", name, elapsed)
}

func NoEcho(...any) {}
