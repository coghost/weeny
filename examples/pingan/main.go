package main

import (
	"github.com/coghost/weeny"
	"github.com/gocolly/redisstorage"
	"github.com/k0kubun/pp/v3"
)

func main() {
	pp.Default.SetExportedOnly(true)

	const (
		home       = `https://talent.pingan.com/recruit/social.html`
		items      = `div.resultList a`
		pagination = `button.btn-next`
	)

	storage := &redisstorage.Storage{
		Address:  "127.0.0.1:6381",
		Password: "",
		DB:       2,
		Prefix:   "pingan",
	}

	c := weeny.NewCrawler()

	c.MustSetStorage(storage)

	c.OnHTML(items, func(e *weeny.HTMLElement) {
		c.Pie(e.Request.VisitElem(
			weeny.OpenInTab(true),
			weeny.OnVisitEnd(c.ClosePage),
		))
	})

	c.OnPaging(pagination, func(e *weeny.SerpElement) {
		c.Echo(e.Request.VisitElem())
	})

	c.EnsureVisit(home)
}
