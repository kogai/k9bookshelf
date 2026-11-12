package content

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	generated "github.com/kogai/k9bookshelf/gqlgenc/client"

	shopify "github.com/bold-commerce/go-shopify/v4"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func deployProducts(shopDomain, shopToken string, contents []Content, bar *mpb.Bar) error {
	gqlClient, ctx := establishGqlClient(shopDomain, shopToken)
	wg := sync.WaitGroup{}
	c := make(chan error)
	for _, content := range contents {
		wg.Add(1)

		go func(handle, html string) {
			defer wg.Done()
			defer bar.Increment()

			productByHandle, err := gqlClient.ProductByHandle(ctx, handle)
			if err != nil {
				c <- err
				return
			}
			if productByHandle.ProductByHandle == nil {
				c <- fmt.Errorf("Product '%s' is not exist on admin", handle)
				return
			}

			res, err := gqlClient.Deploy(
				ctx,
				generated.ProductInput{
					ID:              &productByHandle.ProductByHandle.ID,
					Handle:          &handle,
					DescriptionHTML: &html,
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
		}(content.handle, content.html)
	}
	go func() {
		wg.Wait()
		c <- nil
	}()

	err := <-c
	return err
}

func deployPages(shopDomain, appKey, appSecret string, contents []Content, bar *mpb.Bar) error {
	ctx := context.Background()
	adminClient, err := establishRestClient(shopDomain, appKey, appSecret)
	if err != nil {
		return err
	}
	pages, err := adminClient.Page.List(ctx, nil)
	if err != nil {
		return err
	}

	for _, content := range contents {
		var page *shopify.Page
		for i := range pages {
			if pages[i].Handle == content.handle {
				page = &pages[i]
				break
			}
		}
		if page == nil {
			return fmt.Errorf("page [%s] is not exist", content.handle)
		}

		// NOTE: Fields of shopify.Page are tagged by `omitempty`, so only body is updated.
		_, err = adminClient.Page.Update(ctx, shopify.Page{
			Id:       page.Id,
			BodyHTML: content.html,
		})
		if err != nil {
			return fmt.Errorf("page [%s]: %w", content.handle, err)
		}
		bar.Increment()
	}
	return nil
}

func deployBlogs(shopDomain, appKey, appSecret string, blogs map[string][]Content, bar *mpb.Bar) error {
	adminClient, err := establishRestClient(shopDomain, appKey, appSecret)
	if err != nil {
		return err
	}
	articleResource := NewArticleResource(adminClient)

	currentBlogs, err := adminClient.Blog.List(context.Background(), nil)
	if err != nil {
		return err
	}

	// NOTE: Articles are deployed sequentially to stay within the rate limit of the REST Admin API.
	for blogHandle, blogContents := range blogs {
		var blog *shopify.Blog
		for i := range currentBlogs {
			if currentBlogs[i].Handle == blogHandle {
				blog = &currentBlogs[i]
				break
			}
		}
		if blog == nil {
			blog, err = createBlog(adminClient, blogHandle, blogHandle)
			if err != nil {
				return fmt.Errorf("failed to create blog [%s]: %w", blogHandle, err)
			}
		}

		articles, err := articleResource.List(blog.Id)
		if err != nil {
			return err
		}
		existing := map[string]Article{}
		for _, a := range articles.Articles {
			existing[a.Handle] = a
		}

		for _, content := range blogContents {
			publishedAt, err := content.meta.publishedAt()
			if err != nil {
				return fmt.Errorf("blog article [%s/%s]: %w", blogHandle, content.handle, err)
			}
			article := Article{
				BlogID:      blog.Id,
				BodyHTML:    content.html,
				PublishedAt: publishedAt,
				Tags:        content.meta.tags(),
			}
			if content.meta != nil {
				article.Title = content.meta.Title
			}

			if a, ok := existing[content.handle]; ok {
				article.ID = a.ID
				_, err = articleResource.Put(article)
			} else {
				if article.Title == "" {
					return fmt.Errorf("blog article [%s/%s] is not exist. Add `title` to front matter to create it", blogHandle, content.handle)
				}
				article.Handle = content.handle
				_, err = articleResource.Post(article)
			}
			if err != nil {
				return fmt.Errorf("blog article [%s/%s]: %w", blogHandle, content.handle, err)
			}
			bar.Increment()
		}
	}
	return nil
}

func filesToContents(inputDir string, files []os.FileInfo) ([]Content, error) {
	contents := []Content{}
	for _, file := range files {
		filename := file.Name()
		if file.IsDir() || filepath.Ext(filename) != ".md" {
			continue
		}
		handle := filename[0 : len(filename)-len(filepath.Ext(filename))]
		md, err := ioutil.ReadFile(path.Join(inputDir, filename))
		if err != nil {
			return nil, err
		}
		meta, body, err := splitFrontMatter(string(md))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path.Join(inputDir, filename), err)
		}
		html, err := markdownToHTML(body)
		if err != nil {
			return nil, err
		}
		contents = append(contents, Content{
			handle: handle,
			html:   html,
			meta:   meta,
		})
	}
	return contents, nil
}

type tmpIterable struct {
	f                func(bar *mpb.Bar) error
	numberOfContents int
}

// Deploy uploads contents to store
func Deploy(shopDomain, appKey, appSecret, shopToken, input string) error {
	rawProducts, err := ioutil.ReadDir(path.Join(input, "products"))
	if err != nil {
		return err
	}
	products, err := filesToContents(path.Join(input, "products"), rawProducts)
	if err != nil {
		return err
	}

	rawPages, err := ioutil.ReadDir(path.Join(input, "pages"))
	if err != nil {
		return err
	}
	pages, err := filesToContents(path.Join(input, "pages"), rawPages)

	rawBlogs, err := ioutil.ReadDir(path.Join(input, "blogs"))
	if err != nil {
		return err
	}

	blogs := map[string][]Content{}
	numberOfBlogs := 0
	for _, b := range rawBlogs {
		if b.IsDir() {
			files, err := ioutil.ReadDir(path.Join(input, "blogs", b.Name()))
			if err != nil {
				return err
			}
			contents, err := filesToContents(path.Join(input, "blogs", b.Name()), files)
			if err != nil {
				return err
			}
			numberOfBlogs += len(contents)
			blogs[b.Name()] = contents
		}
	}

	wg := sync.WaitGroup{}
	p := mpb.New(mpb.WithWaitGroup(&wg))
	c := make(chan error)
	for name, _f := range map[string]tmpIterable{
		"products": {
			f: func(bar *mpb.Bar) error {
				return deployProducts(shopDomain, shopToken, products, bar)
			},
			numberOfContents: len(products),
		},
		"pages": {
			f: func(bar *mpb.Bar) error {
				return deployPages(shopDomain, appKey, appSecret, pages, bar)
			},
			numberOfContents: len(pages),
		},
		"blogs": {
			f: func(bar *mpb.Bar) error {
				return deployBlogs(shopDomain, appKey, appSecret, blogs, bar)
			},
			numberOfContents: numberOfBlogs,
		},
	} {
		wg.Add(1)
		bar := p.AddBar(int64(_f.numberOfContents),
			mpb.PrependDecorators(
				decor.Name(path.Join(input, name)),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.AverageETA(decor.ET_STYLE_GO), "done",
				),
			),
		)
		go func(f func(bar *mpb.Bar) error) {
			defer wg.Done()
			if err = f(bar); err != nil {
				c <- err
				return
			}
		}(_f.f)
	}
	go func() {
		p.Wait()
		c <- nil
	}()

	err = <-c
	return err
}
