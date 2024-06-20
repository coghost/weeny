package main

import (
	"github.com/coghost/weeny"
)

func main() {
	c := weeny.NewCrawler(
		weeny.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
		weeny.MaxDepth(3),
	)

	items := `a[href*="wiki"]`

	c.OnHTML(items, func(e *weeny.HTMLElement) {
		c.Pie(e.Request.VisitElem(
			weeny.OnVisitEnd(c.GoBack),
		))
	})

	c.EnsureVisit("https://hackerspaces.org/")
}
