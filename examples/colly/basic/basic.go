package main

import (
	"fmt"

	"github.com/coghost/weeny"
	"github.com/k0kubun/pp/v3"
)

func main() {
	pp.Default.SetExportedOnly(true)

	c := weeny.NewCrawler(
		weeny.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)

	c.OnHTML("a[href]", func(e *weeny.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		e.Request.Visit(link)
	})

	c.Visit("https://hackerspaces.org/")
}
