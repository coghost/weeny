package weeny

import (
	"hash/fnv"
	"io"
	"net/url"
	"strings"
	"unicode/utf8"

	whatwgUrl "github.com/nlnwa/whatwg-url/url"
)

var urlParser = whatwgUrl.NewParser(whatwgUrl.WithPercentEncodeSinglePercentSign())

func normalizeURL(u string) string {
	parsed, err := urlParser.Parse(u)
	if err != nil {
		return u
	}

	return parsed.String()
}

func requestHash(url string, body io.Reader) uint64 {
	h64 := fnv.New64a()
	// reparse the url to fix ambiguities such as
	// "http://example.com" vs "http://example.com/"
	io.WriteString(h64, normalizeURL(url))

	if body != nil {
		io.Copy(h64, body)
	}

	return h64.Sum64()
}

func ParseURL(uri string) (*url.URL, error) {
	u, err := urlParser.Parse(uri)
	if err != nil {
		return nil, err
	}

	u2, err := url.Parse(u.Href(false))
	if err != nil {
		return nil, err
	}

	return u2, err
}

func TruncateString(s string, max int) string {
	if max <= 0 {
		return ""
	}

	if utf8.RuneCountInString(s) < max {
		return s
	}

	return string([]rune(s)[:max])
}

func AbsoluteURL(uri string, baseURL string) (string, error) {
	if strings.HasPrefix(uri, "#") {
		return "", errURLInvalid(uri)
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	ref, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	u := base.ResolveReference(ref)

	return u.String(), nil
}
