package weeny

import (
	"net/url"
)

func (c *Crawler) parseAndCheck(depth int, opts ...VisitOptionFunc) (*url.URL, error) {
	opt := &VisitOptions{}
	bindVisitOptions(opt, opts...)

	if opt.url == "" && opt.elem != nil {
		return nil, nil
	}

	parsedWhatwgURL, err := urlParser.Parse(opt.url)
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(parsedWhatwgURL.Href(false))
	if err != nil {
		return nil, err
	}

	// always allow first visit.
	if c.requestCount == 0 {
		return parsedURL, nil
	}

	// always allow when paination.
	if opt.onType == onTypePaging {
		return nil, nil
	}

	if err := c.requestCheck(parsedURL, depth); err != nil {
		return nil, err
	}

	return parsedURL, nil
}

func (c *Crawler) requestCheck(parsedURL *url.URL, depth int) error {
	if c.maxDepth > 0 && c.maxDepth < depth {
		return ErrMaxDepth
	}

	if err := c.checkVistedStatus(parsedURL); err != nil {
		return err
	}

	if err := c.checkDomains(parsedURL.Hostname()); err != nil {
		return err
	}

	return nil
}

func (c *Crawler) checkDomains(domain string) error {
	if c.allowedDomains == nil || len(c.allowedDomains) == 0 {
		return nil
	}

	for _, d2 := range c.allowedDomains {
		if d2 == domain {
			return nil
		}
	}

	return ErrForbiddenDomain
}
