package weeny

import (
	"github.com/coghost/wee"
)

type HTMLCallback func(e *HTMLElement) error

type htmlCallbackContainer struct {
	Selector string
	Function HTMLCallback
	// DeferFunc func(bot *wee.Bot)
}

type CallbackOptions struct {
	deferFunc func(b *wee.Bot)
}

type CallbackOptionFunc func(o *CallbackOptions)

func bindCallbackOptions(opt *CallbackOptions, opts ...CallbackOptionFunc) {
	for _, f := range opts {
		f(opt)
	}
}

func WithDeferFunc(fn func(b *wee.Bot)) CallbackOptionFunc {
	return func(o *CallbackOptions) {
		o.deferFunc = fn
	}
}
