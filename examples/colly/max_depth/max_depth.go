package main

import (
	"fmt"

	"github.com/coghost/weeny"
)

func main() {
	c := weeny.NewCrawler(
		weeny.MaxDepth(2),
	)

	c.OnHTML("a[href]", func(e *weeny.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		e.Request.Visit(link)
	})

	c.Visit("https://hackerspaces.org/")
}
