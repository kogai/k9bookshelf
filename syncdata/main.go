package main

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
	"time"

	"github.com/Yamashou/gqlgenc/client"
	shopify "github.com/bold-commerce/go-shopify"
	"github.com/gomarkdown/markdown"
	"github.com/mattn/godown"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var apiVersion string = "2020-10"
var rootCmd = &cobra.Command{
	Use:   "datakit",
	Short: "datakit is a content management tool like theme-kit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nothing to do without subcommand.")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Upload contents to store",
	Run: func(cmd *cobra.Command, args []string) {
		err := deploy(cmd.Flag("input").Value.String())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download contents from store",
	Run: func(cmd *cobra.Command, args []string) {
		err := download(cmd.Flag("output").Value.String())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func gqlClient() (*generated.Client, context.Context) {
	authHeader := func(req *http.Request) {
		req.Header.Set("X-Shopify-Access-Token", os.Getenv("MARKDOWN_APP_SECRET"))
	}

	return &generated.Client{
		Client: client.NewClient(http.DefaultClient,
			fmt.Sprintf("https://%s/admin/api/%s/graphql.json", "k9books.myshopify.com", "2020-10"),
			authHeader),
	}, context.Background()
}

func establishRestClient() *shopify.Client {
	app := shopify.App{
		ApiKey:    os.Getenv("MARKDOWN_APP_KEY"),
		ApiSecret: os.Getenv("MARKDOWN_APP_SECRET"),
	}

	return shopify.NewClient(app, "k9books.myshopify.com", os.Getenv("MARKDOWN_APP_SECRET"))
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

func download(output string) error {
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
		articles := Articles{Articles: []Article{}}
		err = restClient.Get(path.Join("admin", "api", apiVersion, "blogs", fmt.Sprint(blog.ID), "articles.json"), &articles, nil)
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

func deploy(input string) error {
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

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	downloadCmd.PersistentFlags().StringP("output", "o", fmt.Sprintf("%s", cwd), "output directory")
	deployCmd.PersistentFlags().StringP("input", "i", fmt.Sprintf("%s", cwd), "input directory")
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(deployCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
