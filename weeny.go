package weeny

import (
	"bytes"
	"minicc/shared/clog"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/coghost/weeny/storage"

	"github.com/PuerkitoBio/goquery"
	"github.com/coghost/wee"
	"github.com/k0kubun/pp/v3"
	"go.uber.org/zap"
)

const (
	_capacity = 4
)

var collectorCounter uint32

type Crawler struct {
	ID  uint32
	Bot *wee.Bot

	maxDepth int

	// allowURLRevisit   bool
	allowedDomains []string
	// disallowedDomains []string

	htmlCallbacks []*htmlCallbackContainer

	store        storage.Storage
	requestCount uint32

	// startTime time.Time
	lock *sync.RWMutex
	wg   *sync.WaitGroup

	logger *zap.Logger
	// _log the base logger, anytime a new search happens, re-set logger to this.
	// WARN: don't call this, use logger instead.
	_log *zap.Logger
}

type CrawlerOption func(*Crawler)

func WithLogger(l *zap.Logger) CrawlerOption {
	return func(c *Crawler) {
		c._log = l
	}
}

func MaxDepth(i int) CrawlerOption {
	return func(c *Crawler) {
		c.maxDepth = i
	}
}

func AllowedDomains(domains ...string) CrawlerOption {
	return func(c *Crawler) {
		c.allowedDomains = domains
	}
}

func NewCrawler(options ...CrawlerOption) *Crawler {
	c := &Crawler{}
	c.init()

	for _, f := range options {
		f(c)
	}

	// post init and bind options
	if c.Bot == nil {
		c.Bot = wee.NewBotDefault()
	}

	return c
}

func (c *Crawler) init() {
	c.ID = atomic.AddUint32(&collectorCounter, 1)
	c._log = clog.MustNewZapLogger()
	c.lock = &sync.RWMutex{}
	c.store = &storage.InMemoryStorage{}
	_ = c.store.Init()
	c.wg = &sync.WaitGroup{}
}

func (c *Crawler) resetLoggerField() {
	c.logger = c._log
}

func (c *Crawler) Visit(url string) error {
	return c.scrape(url, 1, nil)
}

func (c *Crawler) scrape(u string, depth int, ctx *Context) error {
	parsedWhatwgURL, err := urlParser.Parse(u)
	if err != nil {
		return err
	}

	parsedURL, err := url.Parse(parsedWhatwgURL.Href(false))
	if err != nil {
		return err
	}

	if err := c.requestCheck(parsedURL, depth); err != nil {
		return err
	}

	c.wg.Add(1)

	return c.fetch(parsedURL, depth, ctx)
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

func (c *Crawler) fetch(URL *url.URL, depth int, ctx *Context) error {
	defer c.wg.Done()

	request := &Request{
		URL:     URL,
		Depth:   depth,
		crawler: c,
		Ctx:     ctx,
		ID:      atomic.AddUint32(&c.requestCount, 1),
	}

	response, err := c.getPage(URL.String(), request)
	if err != nil {
		return err
	}

	response.Ctx = ctx

	c.Bot.Page().MustWaitDOMStable()

	err = c.handleOnHTML(response)
	if err != nil {
		return err
	}

	return nil
}

func (c *Crawler) getPage(uri string, req *Request) (*Response, error) {
	err := c.Bot.Open(uri)
	if err != nil {
		return nil, err
	}

	return &Response{
		Request: req,
		Body:    []byte(c.Bot.Page().MustHTML()),
	}, nil
}

func (c *Crawler) ClickOpen(serp *SerpElement) {
	// sel, index := request.Selector, request.Index
	pp.Println("clicking", serp.Selector, serp.Index)
	elem, err := c.Bot.Elem(serp.Selector, wee.WithIndex(serp.Index))
	if err != nil {
		panic(err)
	}
	err = c.Bot.ClickElem(elem)
	if err != nil {
		panic(err)
	}
	pp.Println("clicked", serp.Selector, serp.Index)
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

	for _, callback := range c.htmlCallbacks {
		i := 0

		doc.Find(callback.Selector).Each(func(_ int, s *goquery.Selection) {
			for _, n := range s.Nodes {
				e := NewHTMLElementFromSelectionNode(resp, s, n, i)
				i++
				callback.Function(e)
			}
		})
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
