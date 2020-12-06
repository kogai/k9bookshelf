package content

import (
	"bytes"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/gomarkdown/markdown"
	"github.com/tdewolff/minify/v2"
	h "github.com/tdewolff/minify/v2/html"
)

func htmlToMarkdown(html string) (string, error) {
	opt := &md.Options{
		HorizontalRule: "---",
	}

	m := minify.New()
	m.AddFunc("text/html", h.Minify)
	buf := bytes.NewBufferString("")
	if err := m.Minify("text/html", buf, bytes.NewReader([]byte(html))); err != nil {
		return "", err
	}

	converter := md.NewConverter("", true, opt)
	markdownStr, err := converter.ConvertString(buf.String())
	if err != nil {
		return "", err
	}
	return markdownStr + "\n", nil
}

func markdownToHTML(rawmd string) (string, error) {
	html := string(markdown.ToHTML([]byte(rawmd), nil, nil))
	return html, nil
}
