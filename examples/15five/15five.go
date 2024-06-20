package main

import (
	"net/url"
	"strings"

	"github.com/coghost/wee"
	"github.com/coghost/weeny"
	"github.com/gookit/goutil/fsutil"
	"github.com/ungerik/go-dry"
)

const (
	home = "https://www.15five.com/"
)

func main() {
	c := weeny.NewCrawler(
		weeny.AllowedDomains("www.15five.com"),
		weeny.MaxDepth(2),
	)

	c.Bot.DisableImages()

	c.OnHTML("html", func(e *weeny.HTMLElement) {
		raw, err := e.DOM.Html()
		if err != nil {
			panic(err)
		}

		save(e.Request.URL.String(), raw)
	})

	c.OnHTML("a[href]", func(e *weeny.HTMLElement) {
		link := e.Attr("href")
		e.Request.Visit(link)
	})

	c.Visit(home)
}

func save(url string, raw string) {
	name := uniqname(url)
	name = "/tmp/15five/" + name + ".html"
	fsutil.MkParentDir(name)

	err := dry.FileSetString(name, raw)
	if err != nil {
		panic(err)
	}
}

func uniqname(uri string) string {
	pu, err := url.Parse(uri)
	if err != nil {
		panic(err)
	}

	name := uri

	if uri == home {
		return "homepage"
	}

	if pu.Host == "www.15five.com" {
		arr := strings.Split(uri, pu.Host)
		name = arr[1]
	}

	name = strings.TrimPrefix(name, "/")
	name = strings.TrimSuffix(name, "/")

	if name == "" {
		name = "homepage"
	}

	return wee.Filenamify(name)
}
