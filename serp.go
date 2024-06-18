package weeny

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/coghost/wee"
	"github.com/go-rod/rod"
)

type Selector struct {
	Locator string `json:"locator,omitempty" structs:"l,omitempty" yaml:"locator,omitempty"`
	// Index: if the selector got more than 1 elems
	Index int `json:"index,omitempty" structs:"i,omitempty" yaml:"index,omitempty"`
	// Required for popovers/cookies/modal(detail page modal),
	// we want to ensure it has been successfully closed.
	// if not, will return error.
	Required bool `json:"required,omitempty" structs:"required,omitempty" yaml:"required,omitempty"`

	Attr string `json:"attr,omitempty" structs:"attr,omitempty" yaml:"attr,omitempty"`
}

type AttrConfig struct {
	// attr
	Attr      string `json:"attr,omitempty" structs:"attr,omitempty" yaml:"attr,omitempty"`
	AttrSep   string `json:"attr_sep,omitempty" structs:"attr_sep,omitempty" yaml:"attr_sep,omitempty"`
	AttrIndex int    `json:"attr_index,omitempty" structs:"attr_index,omitempty" yaml:"attr_index,omitempty"`
	// AttrChars is used when parse items count, and there exists thousand separator i.e. `"." or ","`
	AttrChars string `json:"attr_chars,omitempty" structs:"attr_chars,omitempty" yaml:"attr_chars,omitempty"`
	// AttrRegex useful if we want to parse datetime from raw
	AttrRegex string `json:"attr_regex,omitempty" structs:"attr_regex,omitempty" yaml:"attr_regex,omitempty"`
	// AttrRefine
	AttrRefine string `json:"attr_refine,omitempty" structs:"attr_refine,omitempty" yaml:"attr_refine,omitempty"`
	AttrJS     string `json:"attr_js,omitempty" structs:"attr_js,omitempty" yaml:"attr_js,omitempty"`
}

type SerpElement struct {
	Request *Request

	DOM *goquery.Selection

	URL      string
	BaseURL  string
	Bot      *wee.Bot
	Selector string
	Index    int
	Element  *rod.Element
}

func NewSerpWithURL(uri string) *SerpElement {
	return &SerpElement{
		URL: uri,
		Request: &Request{
			Depth: 0,
		},
	}
}

func NewSerpElement(req *Request, bot *wee.Bot, sel string, index int) (*SerpElement, error) {
	elem, err := bot.Elem(sel, wee.WithIndex(index))
	if err != nil {
		return nil, err
	}

	serp := &SerpElement{
		Request:  req,
		BaseURL:  bot.CurrentUrl(),
		Bot:      bot,
		Selector: sel,
		Index:    index,
		Element:  elem,
	}

	return serp, nil
}

func (e *SerpElement) Attr(k string) string {
	v, err := e.Element.Attribute(k)
	if err != nil || v == nil {
		return ""
	}

	return *v
}

func (e *SerpElement) Text() string {
	v, err := e.Element.Text()
	if err != nil {
		return ""
	}

	return v
}

func (e *SerpElement) AbsLink(uri string) (string, error) {
	return AbsoluteURL(uri, e.BaseURL)
}

// Link alias of Attr for the first matched of "src/href"
//
//	@return string
func (e *SerpElement) Link() string {
	for _, attr := range []string{"src", "href"} {
		if v := e.Attr(attr); v != "" {
			return v
		}
	}

	return ""
}

func (e *SerpElement) Target() string {
	txt := e.Text()
	link := e.Link()

	if link == "" {
		return txt
	}

	return fmt.Sprintf("%s(%s)", txt, link)
}

func (e *SerpElement) Focus(count int, style string) {
	if count <= 0 {
		return
	}

	visible, _ := e.Element.Visible()
	interactable, _ := e.Element.Interactable()

	if visible && interactable != nil {
		e.Bot.ScrollToElemDirectly(e.Element)
	}

	e.Bot.MarkElems(time.Second*2, e.Element)
}
