package docs

import (
	"bytes"

	"github.com/yuin/goldmark"
)

var md goldmark.Markdown

func init() {
	md = goldmark.New()
}

// RenderMarkdown converts markdown text to HTML.
func RenderMarkdown(source string) string {
	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return source
	}
	return buf.String()
}
