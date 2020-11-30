package content

import (
	"os"
	"path"
	"sync"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

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
			md, err := htmlToMarkdown(descriptionHTML)
			if err != nil {
				c <- err
				return
			}

			_, err = file.WriteString(md)
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
func Download(shopDomain, appKey, appSecret, shopToken, output string) error {
	adminClient, ctx := establishGqlClient(shopDomain, shopToken)
	restClient := establishRestClient(shopDomain, appKey, appSecret)

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
