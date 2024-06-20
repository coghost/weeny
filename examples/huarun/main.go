package main

import (
	"github.com/coghost/weeny"
	"github.com/gocolly/redisstorage"
	"github.com/k0kubun/pp/v3"
)

func main() {
	pp.Default.SetExportedOnly(true)

	const (
		home       = `https://crc.wintalent.cn/wt/CRC/web/index/CompCRCPagerecruit_Social`
		item       = `div.single>p>a`
		pagination = `div.PageShow a[onclick]@@@---@@@下一页`
	)

	storage := &redisstorage.Storage{
		Address:  "127.0.0.1:6381",
		Password: "",
		DB:       2,
		Prefix:   "huarun",
	}

	c := weeny.NewCrawler()

	c.MustSetStorage(storage)

	c.OnHTML(item, func(e *weeny.HTMLElement) {
		c.Pie(e.Request.VisitElem(weeny.OnVisitEnd(c.GoBack)))
	})

	c.OnPaging(pagination, func(e *weeny.SerpElement) {
		c.Echo(e.Request.VisitElem())
	})

	c.EnsureVisit(home)
}
