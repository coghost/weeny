package weeny

import (
	"fmt"
	"os"
	"strings"

	"github.com/coghost/xpretty"
)

const (
	// glyphLink https://fontawesome.com/search?q=link&o=r
	glyphLink = '\uf0c1'
	// glyphElem https://fontawesome.com/search?q=code&o=r
	glyphElem = '\uf121'
	// glyphBug https://fontawesome.com/search?q=bug&o=r
	glyphBug = '\uf188'
)

// echoEachRequest log every request.
func (c *Crawler) echoEachRequest(req *Request) {
	if c.debugRequest {
		fmt.Fprintln(os.Stdout, req.debugString())
	}
}

// echoEachStep log every step.
func (c *Crawler) echoEachStep(format string, args ...any) {
	if c.debugStep {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}

		format = "%s %s " + format
		b := []any{c, xpretty.Yellowf(string(glyphBug))}
		b = append(b, args...)
		fmt.Fprintf(os.Stdout, format, b...)
	}
}
