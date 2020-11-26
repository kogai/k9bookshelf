package syncdata

import (
	"fmt"
	"path"
	"time"

	shopify "github.com/bold-commerce/go-shopify"
)

// Article is documented at https://shopify.dev/docs/admin-api/rest/reference/online-store/article
type Article struct {
	ID                int64      `json:"id"`
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
	err := a.client.Get(path.Join("admin", "api", apiVersion, "blogs", fmt.Sprint(blogID), "articles.json"), &articles, nil)
	if err != nil {
		return nil, err
	}
	return &articles, nil
}

// Put update article
func (a *ArticleResource) Put(article Article) (*Article, error) {
	var response Article
	err := a.client.Put(path.Join("admin", "api", apiVersion, "blogs", fmt.Sprint(article.BlogID), "articles", fmt.Sprintf("%d.json", article.ID)), ArticlePayload{Article: article}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
