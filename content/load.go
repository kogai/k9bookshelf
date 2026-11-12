package content

import (
	"io/ioutil"
	"path"
)

// LocalContents are contents read from a directory like `contents/`.
type LocalContents struct {
	products []Content
	pages    []Content
	// blogs is keyed by handle of blog.
	blogs map[string][]Content
}

// loadContents reads `products`, `pages` and `blogs/<blog-handle>` under input.
func loadContents(input string) (*LocalContents, error) {
	readDir := func(dir string) ([]Content, error) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		return filesToContents(dir, files)
	}

	products, err := readDir(path.Join(input, "products"))
	if err != nil {
		return nil, err
	}
	pages, err := readDir(path.Join(input, "pages"))
	if err != nil {
		return nil, err
	}

	rawBlogs, err := ioutil.ReadDir(path.Join(input, "blogs"))
	if err != nil {
		return nil, err
	}
	blogs := map[string][]Content{}
	for _, b := range rawBlogs {
		if !b.IsDir() {
			continue
		}
		contents, err := readDir(path.Join(input, "blogs", b.Name()))
		if err != nil {
			return nil, err
		}
		blogs[b.Name()] = contents
	}

	return &LocalContents{products: products, pages: pages, blogs: blogs}, nil
}
