package weeny

import (
	"fmt"
	"net/url"
	"time"

	"github.com/avast/retry-go"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/gookit/goutil/sysutil"
)

func (c *Crawler) request(req *Request, uri *url.URL, opts ...VisitOptionFunc) error {
	opt := &VisitOptions{}
	bindVisitOptions(opt, opts...)

	req.ByGetURL = opt.elem == nil
	req.Selector = opt.selector
	req.Index = opt.index

	c.echoEachRequest(req)

	if req.ByGetURL {
		if err := c.Bot.Open(uri.String()); err != nil {
			return err
		}
	} else {
		if err := c.openElem(opt.elem, opts...); err != nil {
			return err
		}
	}

	c.waitPageStable()

	return nil
}

func (c *Crawler) openElem(elem *rod.Element, opts ...VisitOptionFunc) error {
	opt := &VisitOptions{}
	bindVisitOptions(opt, opts...)

	if err := c.tryFocusElem(elem); err != nil {
		c.echoEachStep("cannot focus elem: %+v", err)
		return ErrNotInteractable
	}

	if opt.openInTab {
		return c.openInNewTab(elem)
	}

	return c.Bot.ClickElem(elem)
}

// tryFocusElem
func (c *Crawler) tryFocusElem(elem *rod.Element) error {
	return retry3(
		func() error {
			if _, err := elem.Interactable(); err != nil {
				_ = rod.Try(func() {
					_ = c.Bot.ScrollToElemDirectly(elem)
				})
				return err
			}
			return nil
		})
}

func retry3(retryableFunc retry.RetryableFunc) error {
	return retry.Do(
		retryableFunc,
		retry.LastErrorOnly(true),
		retry.Attempts(3), //nolint:mnd
		retry.Delay(time.Second*1))
}

func (c *Crawler) openInNewTab(elem *rod.Element) error {
	pages := c.Bot.Browser().MustPages()

	ctrlKey := input.ControlLeft
	if sysutil.IsMac() {
		ctrlKey = input.MetaLeft
	}

	c.Bot.FocusAndHighlight(elem)

	err := elem.MustKeyActions().Press(ctrlKey).Type(input.Enter).Do()
	if err != nil {
		return fmt.Errorf("cannot do ctrl click: %w", err)
	}

	err = c.Bot.ActivateLatestOpenedPage(pages, 10) //nolint:mnd
	if err != nil {
		return fmt.Errorf("cannot switch opened page: %w", err)
	}

	return nil
}
