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

	c.OnHTML("a[href]", func(e *weeny.HTMLElement) error {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		return e.Request.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit("https://hackerspaces.org/")
}
