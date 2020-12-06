package content

import (
	"bytes"

	md "github.com/JohannesKaufmann/html-to-markdown"
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
		panic(err)
	}

	converter := md.NewConverter("", true, opt)
	markdownStr, err := converter.ConvertString(buf.String())
	if err != nil {
		return "", err
	}
	return markdownStr + "\n", nil
}
