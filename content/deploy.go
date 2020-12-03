package content

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	generated "github.com/kogai/k9bookshelf/gqlgenc/client"

	shopify "github.com/bold-commerce/go-shopify"
	"github.com/gomarkdown/markdown"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
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
	var err error
	adminClient := establishRestClient(shopDomain, appKey, appSecret)
	wg := sync.WaitGroup{}
	c := make(chan error)
	for _, content := range contents {
		wg.Add(1)

		go func(handle, html string) {
			defer wg.Done()
			defer bar.Increment()

			pages, err := adminClient.Page.List(nil)
			if err != nil {
				c <- err
				return
			}
			var page shopify.Page
			for _, p := range pages {
				if p.Handle == handle {
					page = p
					break
				}
			}

			// NOTE: Because Page struct isn't tagged by `omitempty` and Metafields are initialized with nil,
			metafields := []shopify.Metafield{}
			if page.Metafields != nil {
				metafields = page.Metafields
			}
			_, err = adminClient.Page.Update(shopify.Page{
				ID:             page.ID,
				Author:         page.Author,
				Handle:         handle,
				Title:          page.Title,
				CreatedAt:      page.CreatedAt,
				UpdatedAt:      page.UpdatedAt,
				BodyHTML:       html,
				TemplateSuffix: page.TemplateSuffix,
				PublishedAt:    page.PublishedAt,
				ShopID:         page.ShopID,
				Metafields:     metafields,
			})

			if err != nil {
				c <- err
				return
			}
		}(content.handle, content.html)
	}
	go func() {
		wg.Wait()
		c <- nil
	}()

	err = <-c
	return err
}

func deployBlogs(shopDomain, appKey, appSecret string, blogs map[string][]Content, bar *mpb.Bar) error {
	var err error
	adminClient := establishRestClient(shopDomain, appKey, appSecret)
	wg := sync.WaitGroup{}
	c := make(chan error)

	currentBlogs, err := adminClient.Blog.List(nil)
	if err != nil {
		return err
	}

	for _blogCategory, _blogContents := range blogs {
		wg.Add(1)
		var b *shopify.Blog
		go func(blogCategory string, blogContents []Content) {
			defer wg.Done()

			for _, currentBlog := range currentBlogs {
				if currentBlog.Handle == blogCategory {
					b = &currentBlog
					break
				}
			}
			if b == nil {
				c <- fmt.Errorf("blog category [%s] is not exist", blogCategory)
				return
			}

			articles, err := NewArticleResource(adminClient).List(b.ID)
			if err != nil {
				c <- err
				return
			}

			for _, _content := range blogContents {
				wg.Add(1)
				var _article Article
				for _, a := range articles.Articles {
					if _content.handle == a.Handle {
						_article = a
					}
				}

				if _article.ID == 0 {
					c <- fmt.Errorf("blog article [%s] is not exist", _content.handle)
					return
				}

				go func(content Content, article Article) {
					defer wg.Done()
					defer bar.Increment()

					_, err = NewArticleResource(adminClient).Put(Article{
						ID:       article.ID,
						BlogID:   article.BlogID,
						BodyHTML: content.html,
					})

					if err != nil {
						c <- err
						return
					}
				}(_content, _article)
			}
		}(_blogCategory, _blogContents)

	}
	go func() {
		wg.Wait()
		c <- nil
	}()

	err = <-c
	return err
}

func filesToContents(inputDir string, files []os.FileInfo) ([]Content, error) {
	contents := []Content{}
	for _, file := range files {
		filename := file.Name()
		handle := filename[0 : len(filename)-len(filepath.Ext(filename))]
		md, err := ioutil.ReadFile(path.Join(inputDir, filename))
		if err != nil {
			return nil, err
		}
		html := string(markdown.ToHTML(md, nil, nil))
		contents = append(contents, Content{
			handle: handle,
			html:   html,
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
					decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
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
