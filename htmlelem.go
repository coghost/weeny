package weeny

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type HTMLCallback func(e *HTMLElement)

type htmlCallbackContainer struct {
	Selector string
	Function HTMLCallback
}

// HTMLElement is the representation of a HTML tag.
type HTMLElement struct {
	// Name is the name of the tag
	Name       string
	Text       string
	attributes []html.Attribute
	// Request is the request object of the element's HTML document
	Request *Request
	// Response is the Response object of the element's HTML document
	Response *Response
	// DOM is the goquery parsed DOM object of the page. DOM is relative
	// to the current HTMLElement
	DOM *goquery.Selection
	// GlbIndex stores the position of the current element within all the elements matched by an OnHTML callback
	GlbIndex int

	ElemSelector string
	ElemIndex    int
}

// NewHTMLElementFromSelectionNode creates a HTMLElement from a goquery.Selection Node.
func NewHTMLElementFromSelectionNode(resp *Response, gqSel *goquery.Selection, node *html.Node, glbCbIndex int, elemSel string, elemIndex int) *HTMLElement {
	return &HTMLElement{
		Name:       node.Data,
		Request:    resp.Request,
		Response:   resp,
		Text:       gqSel.Text(),
		DOM:        gqSel,
		GlbIndex:   glbCbIndex,
		attributes: node.Attr,
		ElemIndex:  elemIndex,

		ElemSelector: elemSel,
	}
}

// Attr returns the selected attribute of a HTMLElement or empty string
// if no attribute found
func (h *HTMLElement) Attr(k string) string {
	return h.DOM.AttrOr(k, "")
}

// ChildText returns the concatenated and stripped text content of the matching
// elements.
func (h *HTMLElement) ChildText(goquerySelector string) string {
	return strings.TrimSpace(h.DOM.Find(goquerySelector).Text())
}

// ChildTexts returns the stripped text content of all the matching
// elements.
func (h *HTMLElement) ChildTexts(goquerySelector string) []string {
	var res []string

	h.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {
		res = append(res, strings.TrimSpace(s.Text()))
	})
	return res
}

// ChildAttr returns the stripped text content of the first matching
// element's attribute.
func (h *HTMLElement) ChildAttr(goquerySelector, attrName string) string {
	if attr, ok := h.DOM.Find(goquerySelector).Attr(attrName); ok {
		return strings.TrimSpace(attr)
	}
	return ""
}

// ChildAttrs returns the stripped text content of all the matching
// element's attributes.
func (h *HTMLElement) ChildAttrs(goquerySelector, attrName string) []string {
	var res []string

	h.DOM.Find(goquerySelector).Each(func(_ int, s *goquery.Selection) {
		if attr, ok := s.Attr(attrName); ok {
			res = append(res, strings.TrimSpace(attr))
		}
	})
	return res
}
