package weeny

import (
	"fmt"

	"github.com/coghost/wee"
	"github.com/go-rod/rod"
)

type SerpCallback func(e *SerpElement)

type serpCallbackContainer struct {
	Selector string
	Function SerpCallback
}

type SerpElement struct {
	Request *Request

	Bot      *wee.Bot
	Selector string
	Index    int
	Element  *rod.Element
}

func NewSerpElement(req *Request, bot *wee.Bot, sel string, index int) (*SerpElement, error) {
	elem, err := bot.Elem(sel, wee.WithIndex(index))
	if err != nil {
		return nil, err
	}

	serp := &SerpElement{
		Request:  req,
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

// Link alias of Attr for the first matched of "src/href"
func (e *SerpElement) Link(attrs ...string) string {
	posibleAttrs := []string{"src", "href"}
	posibleAttrs = append(posibleAttrs, attrs...)

	for _, attr := range posibleAttrs {
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
