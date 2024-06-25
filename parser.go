package weeny

import (
	"bytes"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/coghost/xparse"
)

type WeenyParser struct {
	*xparse.HTMLParser
}

func NewParser(raw, cfg []byte) *WeenyParser {
	ps := xparse.NewHTMLParser(raw, cfg)
	return &WeenyParser{ps}
}

func (wp *WeenyParser) ParseToStruct(key string, obj any) error {
	xparse.DoParse(wp)
	return wp.DataAsStruct(obj, key)
}

func HTML2Mkd(body []byte, sel string) (*string, error) {
	converter := md.NewConverter("", true, nil)

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	arr := []string{}

	doc.Find(sel).Each(func(i int, s *goquery.Selection) {
		arr = append(arr, converter.Convert(s))
	})

	raw := strings.Join(arr, "\n")
	return &raw, err
}

func SelectionToMd(sel *goquery.Selection) *string {
	converter := md.NewConverter("", true, nil)
	raw := converter.Convert(sel)
	return &raw
}
