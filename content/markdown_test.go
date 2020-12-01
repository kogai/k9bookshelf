package content

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestHtmlToMarkdown(t *testing.T) {
	t.Parallel()
	md, err := htmlToMarkdown(`<p>abc</p>`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `abc
`, md)
}

func TestHtmlToMarkdownList(t *testing.T) {
	t.Parallel()
	md, err := htmlToMarkdown(`<ul><li>abc</li></ul><p>def</p>`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `- abc

def
`, md)
}

func TestHtmlToMarkdownListAndBold(t *testing.T) {
	t.Parallel()
	md, err := htmlToMarkdown(`<ul><li>abc</li></ul><b>def</b>`)
	assert.Equal(t, nil, err)
	assert.Equal(t, `- abc

 **def**
`, md)
}
