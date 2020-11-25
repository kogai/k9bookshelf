package syncdata

import (
	"fmt"
	"io/ioutil"
	"k9bookshelf/generated"
	"os"
	"path"
	"path/filepath"
	"sync"

	shopify "github.com/bold-commerce/go-shopify"
	"github.com/gomarkdown/markdown"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func deployProducts(contents []Content, bar *mpb.Bar) error {
	gqlClient, ctx := establishGqlClient()
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

func deployPages(contents []Content, bar *mpb.Bar) error {
	var err error
	adminClient := establishRestClient()
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
	f        func(contents []Content, bar *mpb.Bar) error
	contents []Content
}

// Deploy uploads contents to store
func Deploy(input string) error {
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

	wg := sync.WaitGroup{}
	p := mpb.New(mpb.WithWaitGroup(&wg))

	c := make(chan error)
	for name, _f := range map[string]tmpIterable{
		"products": {
			f: func(contents []Content, bar *mpb.Bar) error {
				return deployProducts(products, bar)
			},
			contents: products,
		},
		"pages": {
			f: func(contents []Content, bar *mpb.Bar) error {
				return deployPages(pages, bar)
			},
			contents: pages,
		},
	} {
		wg.Add(1)
		bar := p.AddBar(int64(len(_f.contents)),
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
		go func(f func(contents []Content, bar *mpb.Bar) error) {
			defer wg.Done()
			if err = f(pages, bar); err != nil {
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
