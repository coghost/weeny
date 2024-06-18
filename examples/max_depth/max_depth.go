package main

import (
	"fmt"

	"github.com/coghost/weeny"
)

func main() {
	c := weeny.NewCrawler(
		weeny.MaxDepth(1),
	)

	c.OnHTML("a[href]", func(e *weeny.HTMLElement) error {
		link := e.Attr("href")
		fmt.Println(link)

		return e.Request.Visit(link)
	})

	c.Visit("https://en.wikipedia.org/")
}
