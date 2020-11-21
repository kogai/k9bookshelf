package main

import (
	"context"
	"fmt"
	"k9bookshelf/generated"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

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
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// TODO: Use goroutine
	var downloadGroup sync.WaitGroup

	for _, edge := range res.Products.Edges {
		downloadGroup.Add(1)

		go func(handle, descriptionHTML string) {
			defer downloadGroup.Done()

			file, err := os.Create(path.Join(cwd, "products", handle+".md"))
			if err != nil {
				panic(err)
			}
			err = godown.Convert(file, strings.NewReader(descriptionHTML), nil)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Done: %s.md\n", handle)
		}(edge.Node.Handle, edge.Node.DescriptionHTML)
	}

	downloadGroup.Wait()
}
