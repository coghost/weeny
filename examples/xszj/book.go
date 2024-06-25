package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/coghost/toolbox/pathlib"
	"github.com/coghost/wee"
	"github.com/coghost/weeny"
	"github.com/coghost/weeny/storage"
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

func main() {
	const (
		home  = `https://xszj.org/b/93877/cs/1`
		items = `a[href^="/b/93877/c/"]`

		chapContent = `article.box_con`
	)

	store := storage.MustNewRedisStorage("127.0.0.1:6381", "", 2, weeny.HostFromURL(home))

	c := weeny.NewCrawlerMuted(
		weeny.MaxDepth(2),
		weeny.Headless(false),
	)

	c.MustSetStorage(store)

	c.OnResponse(func(r *weeny.Response) {
		if r.Request.Depth == 1 {
			return
		}

		var chap *Chap
		weeny.NewParser(r.Body, []byte(cfg)).ParseToStruct("chap", &chap)
		if chap.Title == "" {
			return
		}

		rawPth := fmt.Sprintf("/tmp/%s.html", strings.TrimPrefix(r.Request.URL.Path, "/"))
		newFile := pathlib.Path(rawPth)
		dry.PanicIfErr(newFile.MkParentDir())
		newFile.MustSetString(chap.Title + "\n" + strings.Join(chap.Content, "\n"))

		b := wee.Confirm("load next chap?")
		if !b {
			os.Exit(0)
		}
	})

	c.OnHTML(items, func(e *weeny.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
		c.Bot.DumpCookies()
		cookies := c.Bot.Page().MustCookies()

		var httpCks []*http.Cookie
		for _, ck := range cookies {
			hc := &http.Cookie{
				Name:     ck.Name,
				Value:    ck.Value,
				Path:     ck.Path,
				Domain:   ck.Domain,
				Expires:  ck.Expires.Time(),
				Secure:   ck.Secure,
				HttpOnly: ck.HTTPOnly,
			}

			httpCks = append(httpCks, hc)
		}

		dry.FileSetString("/tmp/xszj.cookie", storage.StringifyCookies(httpCks))
	})

	c.EnsureVisit(home)
}
