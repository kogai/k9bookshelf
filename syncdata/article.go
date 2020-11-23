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
	Title             string     `json:"title"`
	CreatedAt         *time.Time `json:"created_at"`
	BodyHTML          string     `json:"body_html"`
	BlogID            int64      `json:"blog_id"`
	Author            string     `json:"author"`
	UserID            int64      `json:"user_id"`
	PublishedAt       *time.Time `json:"published_at"`
	UpdatedAt         *time.Time `json:"updated_at"`
	SummaryHTML       *string    `json:"summary_html"`
	TemplateSuffix    *string    `json:"template_suffix"`
	Handle            string     `json:"handle"`
	Tags              string     `json:"tags"`
	AdminGraphqlAPIID string     `json:"admin_graphql_api_id"`
}

// Articles is not documented yet.
type Articles struct {
	Articles []Article `json:"articles"`
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
