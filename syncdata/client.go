package syncdata

import (
	"context"
	"fmt"
	"k9bookshelf/syncdata/generated"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
	shopify "github.com/bold-commerce/go-shopify"
)

func establishGqlClient() (*generated.Client, context.Context) {
	authHeader := func(req *http.Request) {
		req.Header.Set("X-Shopify-Access-Token", appSecret)
	}

	return &generated.Client{
		Client: client.NewClient(http.DefaultClient,
			fmt.Sprintf("https://%s/admin/api/%s/graphql.json", shopDomain, apiVersion),
			authHeader),
	}, context.Background()
}

func establishRestClient() *shopify.Client {
	app := shopify.App{
		ApiKey:    appKey,
		ApiSecret: appSecret,
	}

	return shopify.NewClient(app, shopDomain, appSecret, shopify.WithVersion(apiVersion))
}

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
