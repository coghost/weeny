package weeny

import (
	"hash/fnv"
	"io"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/coghost/wee"
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
	_, _ = io.WriteString(h64, normalizeURL(url))

	if body != nil {
		_, _ = io.Copy(h64, body)
	}

	return h64.Sum64()
}

func requestNamify(uu *url.URL) string {
	st := ShortenURL(uu)
	st = strings.TrimPrefix(st, "/")
	return wee.Filenamify(st)
}

// HostFromURL extract `Host` from url
func HostFromURL(uri string) string {
	uu, err := ParseURL(uri)
	if err != nil {
		return uri
	}
	return uu.Host
}

// URL2Str return url.String() or "" if url is nil.
func URL2Str(uu *url.URL) string {
	if uu == nil {
		return ""
	}

	return uu.String()
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

func ShortenURL(uu *url.URL) string {
	if uu.RawQuery != "" {
		return uu.RawQuery
	}

	return uu.Path
}

func TruncateString(str string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	if utf8.RuneCountInString(str) < maxLen {
		return str
	}

	return string([]rune(str)[:maxLen])
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
