package content

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"sync"
	"testing"

	shopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/vbauerster/mpb/v8"
	"gopkg.in/go-playground/assert.v1"
)

type recordedRequest struct {
	Method string
	Path   string
	Body   map[string]interface{}
}

// fakeShop is a minimal fake of the REST Admin API for blogs and articles.
type fakeShop struct {
	mu       sync.Mutex
	blogs    []map[string]interface{}
	articles []map[string]interface{}
	requests []recordedRequest
}

func (f *fakeShop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	var body map[string]interface{}
	if b, _ := io.ReadAll(r.Body); len(b) > 0 {
		_ = json.Unmarshal(b, &body)
	}
	f.requests = append(f.requests, recordedRequest{Method: r.Method, Path: r.URL.Path, Body: body})

	write := func(v interface{}) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	prefix := "/admin/api/" + apiVersion
	switch {
	case r.Method == http.MethodGet && r.URL.Path == prefix+"/blogs.json":
		write(map[string]interface{}{"blogs": f.blogs})
	case r.Method == http.MethodPost && r.URL.Path == prefix+"/blogs.json":
		blog := map[string]interface{}{"id": 2}
		for k, v := range body["blog"].(map[string]interface{}) {
			blog[k] = v
		}
		f.blogs = append(f.blogs, blog)
		write(map[string]interface{}{"blog": blog})
	case r.Method == http.MethodGet && r.URL.Path == prefix+"/blogs/2/articles.json":
		write(map[string]interface{}{"articles": f.articles})
	case r.Method == http.MethodPost && r.URL.Path == prefix+"/blogs/2/articles.json":
		write(map[string]interface{}{"article": body["article"]})
	case r.Method == http.MethodPut && r.URL.Path == prefix+"/blogs/2/articles/10.json":
		write(map[string]interface{}{"article": body["article"]})
	default:
		http.Error(w, "unexpected request", http.StatusNotFound)
	}
}

func withFakeShop(t *testing.T, shop http.Handler) string {
	server := httptest.NewServer(shop)
	t.Cleanup(server.Close)
	u, _ := url.Parse(server.URL)

	// NOTE: Rewrite requests for https://<domain> to the fake server.
	transport := &http.Transport{}
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host
		return transport.RoundTrip(r)
	})}
	restClientOptions = []shopify.Option{shopify.WithHTTPClient(client), shopify.WithRetry(0)}
	gqlHTTPClient = client
	t.Cleanup(func() {
		restClientOptions = nil
		gqlHTTPClient = http.DefaultClient
	})
	return "example.myshopify.com"
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestDeployBlogsCreatesBlogAndArticles(t *testing.T) {
	shop := &fakeShop{}
	domain := withFakeShop(t, shop)

	contents := []Content{
		{handle: "new-article", html: "<p>new</p>", meta: &FrontMatter{Title: "New", PublishedAt: "2015-03-30T00:00:00+09:00", Tags: []string{"a", "b"}}},
	}
	err := deployBlogs(domain, "key", "secret", map[string][]Content{"devlog": contents}, mpb.New().AddBar(1))
	assert.Equal(t, err, nil)

	var paths []string
	for _, r := range shop.requests {
		paths = append(paths, r.Method+" "+r.Path)
	}
	assert.Equal(t, paths, []string{
		"GET /admin/api/2020-10/blogs.json",
		"POST /admin/api/2020-10/blogs.json",
		"GET /admin/api/2020-10/blogs/2/articles.json",
		"POST /admin/api/2020-10/blogs/2/articles.json",
	})

	// NOTE: Only title and handle are sent to create blog.
	blog := shop.requests[1].Body["blog"].(map[string]interface{})
	keys := []string{}
	for k := range blog {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	assert.Equal(t, keys, []string{"handle", "title"})
	assert.Equal(t, blog["handle"], "devlog")

	article := shop.requests[3].Body["article"].(map[string]interface{})
	assert.Equal(t, article["handle"], "new-article")
	assert.Equal(t, article["title"], "New")
	assert.Equal(t, article["body_html"], "<p>new</p>")
	assert.Equal(t, article["tags"], "a, b")
	assert.Equal(t, article["published_at"], "2015-03-30T00:00:00+09:00")
}

func TestDeployBlogsUpdatesExistingArticle(t *testing.T) {
	shop := &fakeShop{
		blogs:    []map[string]interface{}{{"id": 2, "handle": "devlog", "title": "devlog"}},
		articles: []map[string]interface{}{{"id": 10, "blog_id": 2, "handle": "existing"}},
	}
	domain := withFakeShop(t, shop)

	contents := []Content{{handle: "existing", html: "<p>updated</p>"}}
	err := deployBlogs(domain, "key", "secret", map[string][]Content{"devlog": contents}, mpb.New().AddBar(1))
	assert.Equal(t, err, nil)

	last := shop.requests[len(shop.requests)-1]
	assert.Equal(t, last.Method+" "+last.Path, "PUT /admin/api/2020-10/blogs/2/articles/10.json")
	article := last.Body["article"].(map[string]interface{})
	assert.Equal(t, article["body_html"], "<p>updated</p>")
	// NOTE: Fields absent from front matter must not be overwritten.
	_, hasTitle := article["title"]
	assert.Equal(t, hasTitle, false)
}

func TestDeployBlogsRequiresTitleToCreateArticle(t *testing.T) {
	shop := &fakeShop{blogs: []map[string]interface{}{{"id": 2, "handle": "devlog"}}}
	domain := withFakeShop(t, shop)

	contents := []Content{{handle: "untitled", html: "<p></p>"}}
	err := deployBlogs(domain, "key", "secret", map[string][]Content{"devlog": contents}, mpb.New().AddBar(1))
	assert.NotEqual(t, err, nil)
}
