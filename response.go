package weeny

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/coghost/toolbox/pathlib"
)

// ResponseCallback is a type alias for OnResponse callback functions
type ResponseCallback func(*Response)

type Response struct {
	Request *Request
	Body    []byte
	Doc     *goquery.Document
	// Ctx is a context between a Request and a Response
	Ctx *Context
}

func (r *Response) Mkd(sel string) (*string, error) {
	return HTML2Mkd(r.Body, sel)
}

// Save writes response body to disk
func (r *Response) Save(fileName string) error {
	return pathlib.Path(fileName).SetBytes(r.Body)
}
