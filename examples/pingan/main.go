package main

import (
	"github.com/coghost/weeny"
	"github.com/coghost/weeny/storage"
	"github.com/k0kubun/pp/v3"
)

func main() {
	pp.Default.SetExportedOnly(true)

	const (
		home       = `https://talent.pingan.com/recruit/social.html`
		items      = `div.resultList a`
		pagination = `button.btn-next`
	)

	c := weeny.NewCrawler()
	c.MustSetStorage(storage.MustNewRedisStorage("127.0.0.1:6381", "", 2, "pingan"))

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
