package content

import (
	"context"
	"fmt"
	"net/http"

	generated "github.com/kogai/k9bookshelf/gqlgenc/client"

	"github.com/Yamashou/gqlgenc/client"
	shopify "github.com/bold-commerce/go-shopify/v4"
)

// gqlHTTPClient is used by GraphQL Admin API client. It's replaced in tests.
var gqlHTTPClient = http.DefaultClient

func establishGqlClient(shopDomain, shopToken string) (*generated.Client, context.Context) {
	authHeader := func(req *http.Request) {
		req.Header.Set("X-Shopify-Access-Token", shopToken)
	}

	return &generated.Client{
		Client: client.NewClient(gqlHTTPClient,
			fmt.Sprintf("https://%s/admin/api/%s/graphql.json", shopDomain, apiVersion),
			authHeader),
	}, context.Background()
}

func establishRestClient(shopDomain, appKey, appSecret string) (*shopify.Client, error) {
	app := shopify.App{
		ApiKey:    appKey,
		ApiSecret: appSecret,
	}

	// NOTE: Retry is applied when the REST Admin API responds 429 (rate limited) or 503.
	opts := append([]shopify.Option{shopify.WithVersion(apiVersion), shopify.WithRetry(3)}, restClientOptions...)
	return shopify.NewClient(app, shopDomain, appSecret, opts...)
}

// restClientOptions are appended to options of REST Admin API client. It's replaced in tests.
var restClientOptions []shopify.Option

func fetchProducts(ctx context.Context, adminClient *generated.Client) (*generated.Products, error) {
	var cursor *string
	var res *generated.Products

	for {
		tmpRes, err := adminClient.Products(ctx, 10, cursor)
		if err != nil {
			return nil, err
		}
		if res == nil {
			res = tmpRes
		} else {
			res.Products.Edges = append(res.Products.Edges, tmpRes.Products.Edges...)
		}

		if !tmpRes.Products.PageInfo.HasNextPage {
			break
		} else {
			last := tmpRes.Products.Edges[len(tmpRes.Products.Edges)-1]
			cursor = &last.Cursor
		}
	}
	return res, nil
}
