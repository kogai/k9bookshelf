package content

import (
	"fmt"
	"path"
	"time"

	shopify "github.com/bold-commerce/go-shopify"
)

// Article is documented at https://shopify.dev/docs/admin-api/rest/reference/online-store/article
type Article struct {
	ID                int64      `json:"id,omitempty"`
	Title             string     `json:"title,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty"`
	BodyHTML          string     `json:"body_html,omitempty"`
	BlogID            int64      `json:"blog_id,omitempty"`
	Author            string     `json:"author,omitempty"`
	UserID            int64      `json:"user_id,omitempty"`
	PublishedAt       *time.Time `json:"published_at,omitempty"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	SummaryHTML       *string    `json:"summary_html,omitempty"`
	TemplateSuffix    *string    `json:"template_suffix,omitempty"`
	Handle            string     `json:"handle,omitempty"`
	Tags              string     `json:"tags,omitempty"`
	AdminGraphqlAPIID string     `json:"admin_graphql_api_id,omitempty"`
}

// Articles is not documented yet.
type Articles struct {
	Articles []Article `json:"articles"`
}

// ArticlePayload is not documented yet.
type ArticlePayload struct {
	Article Article `json:"article"`
}

// ArticleResource is not documented yet.
type ArticleResource struct {
	client *shopify.Client
}

// NewArticleResource is constructor of ArticleResource
func NewArticleResource(client *shopify.Client) *ArticleResource {
	return &ArticleResource{client}
}

// List returns articles
func (a *ArticleResource) List(blogID int64) (*Articles, error) {
	articles := Articles{Articles: []Article{}}
	// NOTE: 250 is the maximum page size of the REST Admin API.
	options := struct {
		Limit int `url:"limit"`
	}{Limit: 250}
	err := a.client.Get(path.Join("admin", "api", apiVersion, "blogs", fmt.Sprint(blogID), "articles.json"), &articles, options)
	if err != nil {
		return nil, err
	}
	return &articles, nil
}

// Post creates article
func (a *ArticleResource) Post(article Article) (*Article, error) {
	var response ArticlePayload
	p := path.Join("admin", "api", apiVersion, "blogs", fmt.Sprint(article.BlogID), "articles.json")
	err := a.client.Post(p, ArticlePayload{Article: article}, &response)
	if err != nil {
		return nil, err
	}
	return &response.Article, nil
}

// Put update article
func (a *ArticleResource) Put(article Article) (*Article, error) {
	var response Article
	p := path.Join("admin", "api", apiVersion, "blogs", fmt.Sprint(article.BlogID), "articles", fmt.Sprintf("%d.json", article.ID))
	err := a.client.Put(p, ArticlePayload{Article: article}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// NOTE: shopify.Blog isn't tagged by `omitempty`, so creating a blog through it sends empty values of every field
// and the API rejects it (e.g. `commentable: must be in a valid state`).
type newBlog struct {
	Title  string `json:"title"`
	Handle string `json:"handle"`
}

type blogPayload struct {
	Blog shopify.Blog `json:"blog"`
}

// createBlog creates a blog only with its title and handle.
func createBlog(client *shopify.Client, title, handle string) (*shopify.Blog, error) {
	var response blogPayload
	p := path.Join("admin", "api", apiVersion, "blogs.json")
	err := client.Post(p, map[string]newBlog{"blog": {Title: title, Handle: handle}}, &response)
	if err != nil {
		return nil, err
	}
	return &response.Blog, nil
}
