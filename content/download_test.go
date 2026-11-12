package content

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

// fakeStore serves products (GraphQL), pages, blogs and articles (REST) for Download.
func fakeStore(t *testing.T) http.Handler {
	prefix := "/admin/api/" + apiVersion
	write := func(w http.ResponseWriter, v interface{}) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == prefix+"/graphql.json":
			write(w, map[string]interface{}{"data": map[string]interface{}{"products": map[string]interface{}{
				"pageInfo": map[string]interface{}{"hasNextPage": false},
				"edges": []interface{}{map[string]interface{}{
					"cursor": "c1",
					"node":   map[string]interface{}{"handle": "a-book", "descriptionHtml": "<p>A <strong>book</strong></p>"},
				}},
			}}})
		case r.Method == http.MethodGet && r.URL.Path == prefix+"/pages.json":
			write(w, map[string]interface{}{"pages": []interface{}{map[string]interface{}{"id": 1, "handle": "about", "body_html": "<p>about</p>"}}})
		case r.Method == http.MethodGet && r.URL.Path == prefix+"/blogs.json":
			write(w, map[string]interface{}{"blogs": []interface{}{map[string]interface{}{"id": 2, "handle": "devlog"}}})
		case r.Method == http.MethodGet && r.URL.Path == prefix+"/blogs/2/articles.json":
			write(w, map[string]interface{}{"articles": []interface{}{map[string]interface{}{
				"id": 10, "blog_id": 2, "handle": "hello", "title": "Hello", "body_html": "<p>hi</p>",
				"published_at": "2015-03-30T00:00:00+09:00", "tags": "a, b",
			}}})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected request", http.StatusNotFound)
		}
	})
}

func TestDownload(t *testing.T) {
	domain := withFakeShop(t, fakeStore(t))
	out := t.TempDir()

	err := Download(domain, "key", "secret", "token", out)
	assert.Equal(t, err, nil)

	read := func(rel string) string {
		b, err := ioutil.ReadFile(filepath.Join(out, rel))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}
	assert.Equal(t, read("products/a-book.md"), "A **book**\n")
	assert.Equal(t, read("pages/about.md"), "about\n")
	assert.Equal(t, read("blogs/devlog/hello.md"), `---
title: Hello
published_at: "2015-03-30T00:00:00+09:00"
tags:
- a
- b
---

hi
`)

	// NOTE: Downloaded article must be deployable as it is.
	meta, body, err := splitFrontMatter(read("blogs/devlog/hello.md"))
	assert.Equal(t, err, nil)
	assert.Equal(t, meta.Title, "Hello")
	assert.Equal(t, meta.tags(), "a, b")
	assert.Equal(t, body, "hi\n")
}

func TestDownloadReportsWriteError(t *testing.T) {
	domain := withFakeShop(t, fakeStore(t))
	out := t.TempDir()
	// NOTE: A file at the place of products/ directory makes writing products fail.
	if err := ioutil.WriteFile(filepath.Join(out, "products"), nil, 0644); err != nil {
		t.Fatal(err)
	}

	err := Download(domain, "key", "secret", "token", out)
	assert.NotEqual(t, err, nil)
}
