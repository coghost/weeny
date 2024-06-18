package weeny

type Response struct {
	Request *Request
	Body    []byte
	// Ctx is a context between a Request and a Response
	Ctx *Context
}
