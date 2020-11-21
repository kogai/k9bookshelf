package main

import (
	"context"
	"fmt"
	"k9bookshelf/generated"
	"net/http"
	"os"
	"strings"

	"github.com/Yamashou/gqlgenc/client"
	"github.com/mattn/godown"
)

func main() {
	authHeader := func(req *http.Request) {
		req.Header.Set("X-Shopify-Access-Token", os.Getenv("MARKDOWN_APP_SECRET"))
		req.Header.Set("Content-Type", "application/json")
	}

	adminClient := &generated.Client{
		Client: client.NewClient(http.DefaultClient,
			fmt.Sprintf("https://%s/admin/api/%s/graphql.json", "k9books.myshopify.com", "2020-10"),
			authHeader),
	}

	ctx := context.Background()
	res, err := adminClient.Products(ctx, 10)

	if err != nil {
		panic(err)
	}
	for _, edge := range res.Products.Edges {
		fmt.Println("")
		fmt.Println(edge.Node.Handle, edge.Node.Metafield)
		err := godown.Convert(os.Stdout, strings.NewReader(edge.Node.DescriptionHTML), nil)
		if err != nil {
			panic(err)
		}
	}
}
