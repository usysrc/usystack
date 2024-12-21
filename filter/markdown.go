package filter

import (
	"html/template"

	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func MarkdownFilter(input template.HTML) template.HTML {
	p := parser.New()
	doc := p.Parse([]byte(input))
	renderer := mdhtml.NewRenderer(mdhtml.RendererOptions{Flags: mdhtml.CommonFlags})
	return template.HTML(markdown.Render(doc, renderer))
}
