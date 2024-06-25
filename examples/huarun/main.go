package main

import (
	"github.com/coghost/weeny"
	"github.com/coghost/weeny/storage"
	"github.com/k0kubun/pp/v3"
)

func main() {
	pp.Default.SetExportedOnly(true)

	const (
		home       = `https://crc.wintalent.cn/wt/CRC/web/index/CompCRCPagerecruit_Social`
		item       = `div.single>p>a`
		pagination = `div.PageShow a[onclick]@@@---@@@下一页`
	)

	c := weeny.NewCrawler()
	c.MustSetStorage(storage.MustNewRedisStorage("127.0.0.1:6381", "", 2, "haurun"))

	c.OnHTML(item, func(e *weeny.HTMLElement) {
		c.Pie(e.Request.VisitElem(weeny.OnVisitEnd(c.GoBack)))
	})

	c.OnPaging(pagination, func(e *weeny.SerpElement) {
		c.Echo(e.Request.VisitElem())
	})

	c.EnsureVisit(home)
}
