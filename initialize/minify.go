package initialize

import (
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
	"regexp"
)

func NewMinify() *minify.M {
	m := minify.New()

	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.Add("text/html", &html.Minifier{
		KeepDefaultAttrVals:     true,
		KeepComments:            true,
		KeepConditionalComments: true,
		KeepDocumentTags:        true,
		KeepEndTags:             true,
		KeepQuotes:              true,
	})
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFunc("application/ld+json", json.Minify)
	m.AddFunc("application/xml", xml.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	return m
}
