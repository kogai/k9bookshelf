package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
)

// EstablishGqlClient create API client of GraphQL
func EstablishGqlClient(shopDomain, apiVersion, appSecret string) (*Client, context.Context) {
	authHeader := func(req *http.Request) {
		req.Header.Set("X-Shopify-Access-Token", appSecret)
	}

	return &Client{
		Client: client.NewClient(http.DefaultClient,
			fmt.Sprintf("https://%s/admin/api/%s/graphql.json", shopDomain, apiVersion),
			authHeader),
	}, context.Background()
}
