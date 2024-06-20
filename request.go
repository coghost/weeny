package weeny

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/coghost/wee"
	"github.com/coghost/xpretty"
	"github.com/go-rod/rod"
)

type onType string

const (
	onTypeHTML   onType = "html"
	onTypePaging onType = "paging"
)

type ErrHandler func(error) error

type VisitOptions struct {
	selector string
	index    int
	onType   onType

	elem      *rod.Element
	openInTab bool
	url       string

	onVisitEnd ErrHandler
}

type VisitOptionFunc func(o *VisitOptions)

func bindVisitOptions(opt *VisitOptions, opts ...VisitOptionFunc) {
	for _, f := range opts {
		f(opt)
	}
}

func withOnType(t onType) VisitOptionFunc {
	return func(o *VisitOptions) {
		o.onType = t
	}
}

func selector(s string) VisitOptionFunc {
	return func(o *VisitOptions) {
		o.selector = s
	}
}

func index(i int) VisitOptionFunc {
	return func(o *VisitOptions) {
		o.index = i
	}
}

// OpenInTab marks request open a link in new tab.
func OpenInTab(b bool) VisitOptionFunc {
	return func(o *VisitOptions) {
		o.openInTab = b
	}
}

func Elem(elem *rod.Element) VisitOptionFunc {
	return func(o *VisitOptions) {
		o.elem = elem
	}
}

func WithURL(u string) VisitOptionFunc {
	return func(o *VisitOptions) {
		o.url = u
	}
}

func OnVisitEnd(fn ErrHandler) VisitOptionFunc {
	return func(o *VisitOptions) {
		o.onVisitEnd = fn
	}
}

type Request struct {
	URL *url.URL
	// ID is the Unique identifier of the request
	ID      uint32
	crawler *Crawler

	Depth int
	// Ctx is a context between a Request and a Response
	Ctx *Context

	baseURL *url.URL

	ByGetURL bool

	onType   onType
	Selector string
	Index    int
}

func (r *Request) RID() string {
	return fmt.Sprintf("[%d.%2d]", r.Depth, r.ID)
}

func (r *Request) String() string {
	elemInfo := ""
	if r.Selector != "" {
		elemInfo = fmt.Sprintf("%s(%d)", r.Selector, r.Index)
	}

	glyph := xpretty.Greenf(string(glyphLink))
	if !r.ByGetURL {
		glyph = xpretty.Cyanf(string(glyphElem))
	}

	uri := ""
	if r.URL != nil {
		uri = ShortenURL(r.URL)
	} else if r.baseURL != nil {
		uri = ShortenURL(r.baseURL)
	}

	return fmt.Sprintf("%s %s %s %s", r.RID(), glyph, elemInfo, uri)
}

func (r *Request) debugString() string {
	symb := ""
	sybDepth := r.Depth - 1

	if sybDepth > 0 {
		symb = "└──"
	}

	lead := ""
	if sybDepth > 1 {
		lead = strings.Repeat("   ", sybDepth-1)
	}

	return fmt.Sprintf("%s%s%s", lead, symb, r.String())
}

// AbsoluteURL returns with the resolved absolute URL of an URL chunk.
// AbsoluteURL returns empty string if the URL chunk is a fragment or
// could not be parsed
func (r *Request) AbsoluteURL(uri string) string {
	if strings.HasPrefix(uri, "#") {
		return ""
	}

	base := r.URL
	if r.baseURL != nil {
		base = r.baseURL
	}

	absURL, err := urlParser.ParseRef(base.String(), uri)
	if err != nil {
		return ""
	}
	return absURL.Href(false)
}

// Visit continues Collector's collecting job by creating a
// request and preserves the Context of the previous request.
// Visit also calls the previously provided callbacks
func (r *Request) Visit(url string, opts ...VisitOptionFunc) error {
	opts = append(opts, WithURL(url))
	return r.visit(opts...)
}

func (r *Request) VisitElem(opts ...VisitOptionFunc) error {
	opt := &VisitOptions{}
	bindVisitOptions(opt, opts...)

	// auto get elem by selector and index
	if opt.elem == nil {
		elem, err := r.crawler.Bot.Elem(r.Selector, wee.WithIndex(r.Index))
		if err != nil {
			return err
		}

		if elem == nil {
			return ErrNoElemFound
		}

		opt.elem = elem
		opts = append(opts, Elem(elem))
	}

	// try get url
	if opt.elem != nil && opt.url == "" {
		attr := ""

		href, _ := opt.elem.Attribute("href")
		if href != nil {
			attr = *href
		}

		opts = append(opts, WithURL(attr))
	}

	return r.visit(opts...)
}

func (r *Request) visit(opts ...VisitOptionFunc) error {
	opt := &VisitOptions{}
	bindVisitOptions(opt, opts...)

	if u := opt.url; u != "" {
		u = r.AbsoluteURL(u)
		// after update url, re append to overwrite it.
		opts = append(opts, WithURL(u))
	}

	opts = append(opts, selector(r.Selector), index(r.Index), withOnType(r.onType))

	err := r.crawler.scrape(r.Ctx, r.Depth+1, opts...)
	if opt.onVisitEnd == nil {
		return err
	} else {
		return opt.onVisitEnd(err)
	}
}
