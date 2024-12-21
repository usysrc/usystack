package filter

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
)

func MarkdownFilter(input template.HTML) template.HTML {
	p := parser.New()
	doc := p.Parse([]byte(input))
	renderer := mdhtml.NewRenderer(mdhtml.RendererOptions{Flags: mdhtml.CommonFlags})
	content := markdown.Render(doc, renderer)
	policy := bluemonday.UGCPolicy() // User-Generated Content policy.
	sanitized := policy.SanitizeBytes(content)
	return template.HTML(sanitized)
}
