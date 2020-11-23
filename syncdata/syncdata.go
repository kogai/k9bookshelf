package syncdata

import (
	"context"
	"fmt"
	"io/ioutil"
	"k9bookshelf/generated"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Yamashou/gqlgenc/client"
	shopify "github.com/bold-commerce/go-shopify"
	"github.com/gomarkdown/markdown"
	"github.com/mattn/godown"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

const apiVersion string = "2020-10"
const shopDomain string = "k9books.myshopify.com"

var appKey string = os.Getenv("MARKDOWN_APP_KEY")
var appSecret string = os.Getenv("MARKDOWN_APP_SECRET")
var shopToken string = appSecret

func gqlClient() (*generated.Client, context.Context) {
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

// Content is not documented yet.
type Content struct {
	handle string
	html   string
}

// Contents is not documented yet.
type Contents struct {
	kind  string
	items []Content
}

func dowloadContens(output string, contents *[]Content, bar *mpb.Bar) error {
	err := os.MkdirAll(output, os.ModePerm)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for _, content := range *contents {
		wg.Add(1)
		c := make(chan error)

		go func(handle, descriptionHTML string) {
			defer wg.Done()
			defer bar.Increment()

			file, err := os.Create(path.Join(output, handle+".md"))
			if err != nil {
				c <- err
				return
			}
			err = godown.Convert(file, strings.NewReader(descriptionHTML), nil)
			if err != nil {
				c <- err
				return
			}
			c <- nil
		}(content.handle, content.html)
		err = <-c
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

// Download downloads contents from store
func Download(output string) error {
	adminClient, ctx := gqlClient()
	restClient := establishRestClient()

	res, err := fetchProducts(ctx, adminClient)
	if err != nil {
		return err
	}
	products := Contents{kind: "products"}
	for _, product := range res.Products.Edges {
		products.items = append(products.items, Content{
			handle: product.Node.Handle,
			html:   product.Node.DescriptionHTML,
		})
	}

	rawPages, err := restClient.Page.List(nil)
	if err != nil {
		return err
	}
	pages := Contents{kind: "pages"}
	for _, page := range rawPages {
		pages.items = append(pages.items, Content{
			handle: page.Handle,
			html:   page.BodyHTML,
		})
	}

	rawBlogs, err := restClient.Blog.List(nil)
	if err != nil {
		return err
	}
	blogs := []Contents{}
	for _, blog := range rawBlogs {
		contents := Contents{
			kind: path.Join("blogs", blog.Handle),
		}
		articles, err := NewArticleResource(restClient).List(blog.ID)
		if err != nil {
			return err
		}
		for _, article := range articles.Articles {
			contents.items = append(contents.items, Content{
				handle: article.Handle,
				html:   article.BodyHTML,
			})
		}
		blogs = append(blogs, contents)
	}

	var wg sync.WaitGroup
	progress := mpb.New(mpb.WithWaitGroup(&wg))
	for _, cts := range append([]Contents{products, pages}, blogs...) {
		wg.Add(1)
		bar := progress.AddBar(int64(len(cts.items)),
			mpb.PrependDecorators(
				decor.Name(cts.kind),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
				),
			),
		)

		go func(o string, items []Content, b *mpb.Bar) {
			defer wg.Done()
			err = dowloadContens(o, &items, b)
		}(path.Join(output, cts.kind), cts.items, bar)
		if err != nil {
			return err
		}
	}

	progress.Wait()
	return nil
}

// Deploy uploads contents to store
func Deploy(input string) error {
	files, err := ioutil.ReadDir(path.Join(input, "products"))
	if err != nil {
		return err
	}
	adminClient, ctx := gqlClient()
	wg := sync.WaitGroup{}
	p := mpb.New(mpb.WithWaitGroup(&wg))
	bar := p.AddBar(int64(len(files)),
		mpb.PrependDecorators(
			decor.Name(path.Join(input, "products")),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
			),
		),
	)

	for _, file := range files {
		wg.Add(1)
		c := make(chan error)
		filename := file.Name()

		go func(handle, pathToFile string) {
			defer wg.Done()
			defer bar.Increment()

			productByHandle, err := adminClient.ProductByHandle(ctx, handle)
			if err != nil {
				c <- err
				return
			}

			md, err := ioutil.ReadFile(pathToFile)
			if err != nil {
				c <- err
				return
			}
			descriptionHTML := string(markdown.ToHTML(md, nil, nil))

			res, err := adminClient.Deploy(
				ctx,
				generated.ProductInput{
					ID:              &productByHandle.ProductByHandle.ID,
					Handle:          &handle,
					DescriptionHTML: &descriptionHTML,
				},
			)
			if err != nil {
				c <- err
				return
			}
			if len(res.ProductUpdate.UserErrors) > 0 {
				var errorBuf string
				for _, userError := range res.ProductUpdate.UserErrors {
					errorBuf += fmt.Sprintf("'%s': '%s'\n", userError.Field, userError.Message)
				}
				c <- fmt.Errorf("{\n%s}", errorBuf)
				return
			}
			c <- nil
		}(
			filename[0:len(filename)-len(filepath.Ext(filename))],
			path.Join(input, "products", filename),
		)

		err = <-c
		if err != nil {
			return err
		}
	}
	p.Wait()
	return nil
}
