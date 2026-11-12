package content

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSplitFrontMatter(t *testing.T) {
	fm, body, err := splitFrontMatter("---\ntitle: \"タイトル\"\npublished_at: \"2015-03-30T00:00:00+09:00\"\ntags:\n  - a\n  - b\n---\n\nbody\n\n---\n\nafter\n")
	assert.Equal(t, err, nil)
	assert.Equal(t, fm.Title, "タイトル")
	assert.Equal(t, fm.tags(), "a, b")
	assert.Equal(t, body, "body\n\n---\n\nafter\n")
	p, err := fm.publishedAt()
	assert.Equal(t, err, nil)
	assert.Equal(t, p.Unix(), int64(1427641200))
}

func TestSplitFrontMatterWithoutFrontMatter(t *testing.T) {
	fm, body, err := splitFrontMatter("abc\n---\ndef\n")
	assert.Equal(t, err, nil)
	assert.Equal(t, fm == nil, true)
	assert.Equal(t, body, "abc\n---\ndef\n")
}

func TestSplitFrontMatterUnclosed(t *testing.T) {
	_, _, err := splitFrontMatter("---\ntitle: a\n")
	assert.NotEqual(t, err, nil)
}

func TestFrontMatterRoundTrip(t *testing.T) {
	in := &FrontMatter{Title: "t: colon", PublishedAt: "2015-03-30T00:00:00+09:00", Tags: []string{"作ったもの"}}
	rendered, err := renderFrontMatter(in)
	assert.Equal(t, err, nil)
	out, body, err := splitFrontMatter(rendered + "body\n")
	assert.Equal(t, err, nil)
	assert.Equal(t, out, in)
	assert.Equal(t, body, "body\n")
}
