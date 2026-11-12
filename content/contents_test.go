package content

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Run `make golden` (`go test ./content -run TestRenderedContents -update`) to regenerate golden files.
var update = flag.Bool("update", false, "update golden files of rendered contents")

const (
	contentsDir = "../contents"
	goldenDir   = "testdata/rendered"
)

// NOTE: Shopify handles are lowercase. Non-ASCII letters (e.g. Japanese) are allowed as existing handles use them.
var invalidHandle = regexp.MustCompile(`[A-Z\s_./]`)

// TestContentFiles checks that every file under contents/ is deployable.
func TestContentFiles(t *testing.T) {
	err := filepath.Walk(contentsDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(contentsDir, p)
		if info.IsDir() || strings.HasPrefix(rel, "assets"+string(filepath.Separator)) {
			return nil
		}

		parts := strings.Split(filepath.ToSlash(rel), "/")
		kind := parts[0]
		switch {
		case (kind == "products" || kind == "pages") && len(parts) == 2:
		case kind == "blogs" && len(parts) == 3:
		default:
			t.Errorf("%s: unexpected location. Contents must be at products/<handle>.md, pages/<handle>.md or blogs/<blog>/<handle>.md", rel)
			return nil
		}
		if filepath.Ext(p) != ".md" {
			t.Errorf("%s: not a markdown file. It is ignored by deploy", rel)
			return nil
		}
		handle := strings.TrimSuffix(filepath.Base(p), ".md")
		if invalidHandle.MatchString(handle) {
			t.Errorf("%s: handle %q must be lowercase without spaces, underscores or dots", rel, handle)
		}

		raw, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		meta, body, err := splitFrontMatter(string(raw))
		if err != nil {
			t.Errorf("%s: invalid front matter: %v", rel, err)
			return nil
		}
		if strings.TrimSpace(body) == "" {
			t.Errorf("%s: body is empty", rel)
		}
		if meta == nil {
			return nil
		}
		if kind != "blogs" {
			t.Errorf("%s: front matter is supported only by blog articles", rel)
			return nil
		}
		if meta.Title == "" {
			t.Errorf("%s: front matter requires title", rel)
		}
		if _, err := meta.publishedAt(); err != nil {
			t.Errorf("%s: %v", rel, err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestRenderedContents compares HTML rendered from contents/ with golden files,
// so that a change of rendering (e.g. by updating markdown library) shows up as a diff.
func TestRenderedContents(t *testing.T) {
	loaded, err := loadContents(contentsDir)
	if err != nil {
		t.Fatal(err)
	}

	rendered := map[string]string{}
	add := func(dir string, contents []Content) {
		for _, c := range contents {
			rendered[path.Join(dir, c.handle+".html")] = c.html
		}
	}
	add("products", loaded.products)
	add("pages", loaded.pages)
	for blog, contents := range loaded.blogs {
		add(path.Join("blogs", blog), contents)
	}

	if *update {
		if err := os.RemoveAll(goldenDir); err != nil {
			t.Fatal(err)
		}
		for name, html := range rendered {
			p := filepath.Join(goldenDir, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
				t.Fatal(err)
			}
			if err := ioutil.WriteFile(p, []byte(html), 0644); err != nil {
				t.Fatal(err)
			}
		}
		return
	}

	golden := map[string]bool{}
	err = filepath.Walk(goldenDir, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(goldenDir, p)
		golden[filepath.ToSlash(rel)] = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	names := make([]string, 0, len(rendered))
	for name := range rendered {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		want, err := ioutil.ReadFile(filepath.Join(goldenDir, filepath.FromSlash(name)))
		if err != nil {
			t.Errorf("%s: golden file is missing. Run `make golden`", name)
			continue
		}
		if string(want) != rendered[name] {
			t.Errorf("%s: rendered HTML differs from golden file. Run `make golden` and review the diff", name)
		}
		delete(golden, name)
	}
	for name := range golden {
		t.Errorf("%s: golden file has no corresponding content. Run `make golden`", name)
	}
}
