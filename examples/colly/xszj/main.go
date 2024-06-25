package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/coghost/toolbox/pathlib"
	"github.com/coghost/toolbox/sleep"
	"github.com/coghost/weeny"
	"github.com/gocolly/colly"
	"github.com/ungerik/go-dry"
)

type Chap struct {
	Title   string
	Content []string
}

var cfg = `
chap:
  title: h1.bookname
  content:
    _l: div#booktxt p
    _i: ~
`

var bookID = flag.String("id", "", "book id")

func main() {
	flag.Parse()

	if *bookID == "" {
		fmt.Printf("please specify the bookID, -id <BOOKID>\n")
		os.Exit(0)
	}

	var (
		home  = fmt.Sprintf(`https://xszj.org/b/%s/cs/1`, *bookID)
		items = fmt.Sprintf(`a[href^="/b/%s/c/"]`, *bookID)
	)

	c := colly.NewCollector()

	c.OnHTML(items, func(e *colly.HTMLElement) {
		err := e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
		if err == nil {
			fmt.Printf("%s -> %s\n", e.Text, e.Attr("href"))
		}
		sleep.PT10Ms()
	})

	c.OnResponse(func(r *colly.Response) {
		var chap *Chap
		weeny.NewParser(r.Body, []byte(cfg)).ParseToStruct("chap", &chap)
		if chap.Title == "" {
			return
		}

		newFile := pathlib.Path(fmt.Sprintf("/tmp/%s.html", strings.TrimPrefix(r.Request.URL.Path, "/")))
		dry.PanicIfErr(newFile.MkParentDir())
		newFile.MustSetString(chap.Title + "\n" + strings.Join(chap.Content, "\n"))
	})

	dry.PanicIfErr(c.Visit(home))
}
