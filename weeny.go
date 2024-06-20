package weeny

import (
	"bytes"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
	"github.com/coghost/weeny/storage"
	"github.com/coghost/zlog"

	"github.com/coghost/wee"
	"go.uber.org/zap"
)

const (
	_capacity = 4
)

var collectorCounter uint32

type Crawler struct {
	ID  uint32
	Bot *wee.Bot

	maxDepth     int
	requestCount uint32

	debugRequest bool
	debugStep    bool

	stableDiff float64
	headless   bool
	trackTime  bool

	// allowURLRevisit   bool
	allowedDomains []string
	// disallowedDomains []string
	ignoredErrors []error

	htmlCallbacks   []*htmlCallbackContainer
	pagingCallbacks []*serpCallbackContainer

	store storage.Storage

	lock *sync.RWMutex
	wg   *sync.WaitGroup

	logger *zap.Logger
}

type CrawlerOption func(*Crawler)

func WithLogger(l *zap.Logger) CrawlerOption {
	return func(c *Crawler) {
		c.logger = l
	}
}

func DebugEachRequest(b bool) CrawlerOption {
	return func(c *Crawler) {
		c.debugRequest = b
	}
}

func DebugDetailStep(b bool) CrawlerOption {
	return func(c *Crawler) {
		c.debugStep = b
	}
}

func MaxDepth(i int) CrawlerOption {
	return func(c *Crawler) {
		c.maxDepth = i
	}
}

func StableDiff(f float64) CrawlerOption {
	return func(c *Crawler) {
		c.stableDiff = f
	}
}

func Headless(b bool) CrawlerOption {
	return func(c *Crawler) {
		c.headless = b
	}
}

func TrackTime(b bool) CrawlerOption {
	return func(c *Crawler) {
		c.trackTime = b
	}
}

func AllowedDomains(domains ...string) CrawlerOption {
	return func(c *Crawler) {
		c.allowedDomains = domains
	}
}

func IgnoredErrors(errs ...error) CrawlerOption {
	return func(c *Crawler) {
		c.ignoredErrors = errs
	}
}

func NewCrawlerMuted(options ...CrawlerOption) *Crawler {
	c := NewCrawler()

	c.debugRequest = false
	c.debugStep = false

	return c
}

func NewCrawler(options ...CrawlerOption) *Crawler {
	c := &Crawler{}
	c.init()

	for _, f := range options {
		f(c)
	}

	// post init and bind options
	if c.Bot == nil {
		c.Bot = wee.NewBotDefault(wee.Headless(c.headless))
	}

	return c
}

func (c *Crawler) init() {
	c.ID = atomic.AddUint32(&collectorCounter, 1)
	c.logger = zlog.MustNewLoggerDebug()
	c.lock = &sync.RWMutex{}
	c.store = &storage.InMemoryStorage{}
	_ = c.store.Init()
	c.wg = &sync.WaitGroup{}

	c.debugRequest = true
	c.debugStep = true
	c.stableDiff = 0.05
}

// String is the text representation of the crawler.
// It contains useful debug information about the collector's internals
func (c *Crawler) String() string {
	return fmt.Sprintf(
		"<%d.%2d>",
		c.ID,
		atomic.LoadUint32(&c.requestCount),
	)
}

func (c *Crawler) MustSetStorage(s storage.Storage) {
	if err := c.SetStorage(s); err != nil {
		panic(err)
	}
}

func (c *Crawler) SetStorage(s storage.Storage) error {
	if err := s.Init(); err != nil {
		return err
	}

	c.store = s

	return nil
}

func (c *Crawler) EnsureVisit(url string, opts ...VisitOptionFunc) {
	opt := &VisitOptions{onVisitEnd: c.filterErrors}
	bindVisitOptions(opt, opts...)

	if err := c.Visit(url, opts...); err != nil {
		if err = opt.onVisitEnd(err); err != nil {
			c.logger.Error("visit failed", zap.Error(err))
		}
	}
}

func (c *Crawler) Visit(url string, opts ...VisitOptionFunc) error {
	opts = append(opts, WithURL(url))
	return c.scrape(nil, 1, opts...)
}

func (c *Crawler) scrape(ctx *Context, depth int, opts ...VisitOptionFunc) error {
	opt := &VisitOptions{}
	bindVisitOptions(opt, opts...)

	// c.echoEachStep("check(%s)", opt.url)
	parsedURL, err := c.parseAndCheck(depth, opts...)
	if err != nil {
		// c.echoEachStep("parseURL failed: %+v", err)
		return fmt.Errorf("parse or check failed: %w", err)
	}

	c.wg.Add(1)

	return c.fetch(parsedURL, depth, ctx, opts...)
}

func (c *Crawler) fetch(parsedURL *url.URL, depth int, ctx *Context, opts ...VisitOptionFunc) error {
	defer c.wg.Done()

	request := &Request{
		URL:     parsedURL,
		Depth:   depth,
		crawler: c,
		Ctx:     ctx,
		ID:      atomic.AddUint32(&c.requestCount, 1),
	}

	err := c.request(request, parsedURL, opts...)
	if err != nil {
		return err
	}

	response := &Response{
		Request: request,
		Body:    []byte(c.Bot.Page().MustHTML()),
	}

	response.Ctx = ctx

	c.echoEachStep("OnHTML %s", request.RID())
	if err := c.handleOnHTML(response); err != nil {
		return err
	}

	c.echoEachStep("OnPaging %s", request.RID())
	if err := c.handleOnPaging(response); err != nil {
		return err
	}

	return nil
}

func (c *Crawler) checkVistedStatus(parsedURL *url.URL) error {
	uri := parsedURL.String()
	uHash := requestHash(uri, nil)

	visited, err := c.store.IsVisited(uHash)
	if err != nil {
		return err
	}

	if visited {
		return errVisited(uri)
	}

	return c.store.Visited(uHash)
}

// OnHTML registers a function. Function will be executed on every HTML
// element matched by the GoQuery Selector parameter.
// GoQuery Selector is a selector used by https://github.com/PuerkitoBio/goquery
func (c *Crawler) OnHTML(sel string, cbFn HTMLCallback) {
	c.lock.Lock()
	if c.htmlCallbacks == nil {
		c.htmlCallbacks = make([]*htmlCallbackContainer, 0, _capacity)
	}

	c.htmlCallbacks = append(c.htmlCallbacks, &htmlCallbackContainer{
		Selector: sel,
		Function: cbFn,
	})

	c.lock.Unlock()
}

// OnHTMLDetach deregister a function. Function will not be execute after detached
func (c *Crawler) OnHTMLDetach(selector string) {
	c.lock.Lock()

	deleteIdx := -1

	for i, cc := range c.htmlCallbacks {
		if cc.Selector == selector {
			deleteIdx = i
			break
		}
	}

	if deleteIdx != -1 {
		c.htmlCallbacks = append(c.htmlCallbacks[:deleteIdx], c.htmlCallbacks[deleteIdx+1:]...)
		c.logger.Info("detached HTML handler", zap.String("selector", selector))
	}

	c.lock.Unlock()
}

func (c *Crawler) OnPaging(selector string, f SerpCallback) {
	c.lock.Lock()

	c.pagingCallbacks = append(c.pagingCallbacks, &serpCallbackContainer{
		Selector: selector,
		Function: f,
	})

	c.lock.Unlock()
}

func (c *Crawler) handleOnHTML(resp *Response) error {
	if len(c.htmlCallbacks) == 0 {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body))
	if err != nil {
		return err
	}

	if href, found := doc.Find("base[href]").Attr("href"); found {
		u, err := urlParser.ParseRef(resp.Request.URL.String(), href)
		if err == nil {
			baseURL, err := url.Parse(u.Href(false))
			if err == nil {
				resp.Request.baseURL = baseURL
			}
		}
	}

	// try again if href not found by goquery.
	if resp.Request.baseURL == nil {
		baseURL, err := url.Parse(c.Bot.CurrentUrl())
		if err == nil {
			resp.Request.baseURL = baseURL
		}
	}

	for _, callback := range c.htmlCallbacks {
		glbCbIndex := 0

		selection := doc.Find(callback.Selector)
		if len(selection.Nodes) == 0 {
			return ErrNoElemFound
		}

		selection.Each(func(elemIndex int, s *goquery.Selection) {
			for _, node := range s.Nodes {
				resp.Request.Selector = callback.Selector
				resp.Request.Index = elemIndex

				elem := NewHTMLElementFromSelectionNode(resp, s, node, glbCbIndex, callback.Selector, elemIndex)
				glbCbIndex++

				callback.Function(elem)
			}
		})
	}

	return nil
}

func (c *Crawler) handleOnPaging(resp *Response) error {
	resp.Request.Depth = 1
	return c.handleOnSerp(resp, onTypePaging, c.pagingCallbacks)
}

func (c *Crawler) handleOnSerp(resp *Response, ontype onType, callbacks []*serpCallbackContainer) error {
	if len(callbacks) == 0 {
		return nil
	}

	bot := c.Bot

	for _, callback := range callbacks {
		sel := callback.Selector

		// unlike goquery based on HTML crawled,
		// on serp is usually used for dynamical loading pages.
		// so each time, we wait elem shown or quit.
		elems, err := bot.Elems(sel)
		if err != nil {
			return fmt.Errorf("cannot get serp elem %s: %w", sel, err)
		}

		if len(elems) == 0 {
			return ErrNoElemFound
		}

		for index := 0; index < len(elems); index++ {
			resp.Request.Selector = sel
			resp.Request.Index = index
			resp.Request.onType = ontype

			// check total elems, elems may change when we reget all elements.
			// c.logger.Sugar().Debugf("found %s, %d, %d", sel, found, index)
			serpElem, err := NewSerpElement(resp.Request, bot, sel, index)
			if err != nil {
				break
			}

			maxLen := 32
			txt := serpElem.Target()

			if len(txt) > maxLen {
				txt = TruncateString(serpElem.Target(), maxLen) + "..."
			}

			target := fmt.Sprintf("%s | %s=%s", resp.Request.RID(), ontype, txt)
			c.echoEachStep("handle %s %s", ontype, target)

			callback.Function(serpElem)
		}
	}

	return nil
}
