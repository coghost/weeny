package weeny

import (
	"net/url"
	"strings"
)

type Request struct {
	URL *url.URL
	// ID is the Unique identifier of the request
	ID      uint32
	crawler *Crawler

	Depth int
	// Ctx is a context between a Request and a Response
	Ctx *Context

	baseURL *url.URL
}

// AbsoluteURL returns with the resolved absolute URL of an URL chunk.
// AbsoluteURL returns empty string if the URL chunk is a fragment or
// could not be parsed
func (r *Request) AbsoluteURL(u string) string {
	if strings.HasPrefix(u, "#") {
		return ""
	}
	var base *url.URL
	if r.baseURL != nil {
		base = r.baseURL
	} else {
		base = r.URL
	}

	absURL, err := urlParser.ParseRef(base.String(), u)
	if err != nil {
		return ""
	}
	return absURL.Href(false)
}

// Visit continues Collector's collecting job by creating a
// request and preserves the Context of the previous request.
// Visit also calls the previously provided callbacks
func (r *Request) Visit(URL string) error {
	return r.crawler.scrape(r.AbsoluteURL(URL), r.Depth+1, r.Ctx)
}
