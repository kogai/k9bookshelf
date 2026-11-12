package content

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func dowloadContens(output string, contents []Content, bar *mpb.Bar) error {
	if err := os.MkdirAll(output, os.ModePerm); err != nil {
		return err
	}
	for _, content := range contents {
		md, err := htmlToMarkdown(content.html)
		if err != nil {
			return fmt.Errorf("%s: %w", path.Join(output, content.handle), err)
		}
		fm, err := renderFrontMatter(content.meta)
		if err != nil {
			return fmt.Errorf("%s: %w", path.Join(output, content.handle), err)
		}
		if err := os.WriteFile(path.Join(output, content.handle+".md"), []byte(fm+md), 0644); err != nil {
			return err
		}
		bar.Increment()
	}
	return nil
}

// Download downloads contents from store
func Download(shopDomain, appKey, appSecret, shopToken, output string) error {
	adminClient, ctx := establishGqlClient(shopDomain, shopToken)
	restClient, err := establishRestClient(shopDomain, appKey, appSecret)
	if err != nil {
		return err
	}

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

	rawPages, err := restClient.Page.List(context.Background(), nil)
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

	rawBlogs, err := restClient.Blog.List(context.Background(), nil)
	if err != nil {
		return err
	}
	blogs := []Contents{}
	for _, blog := range rawBlogs {
		contents := Contents{
			kind: path.Join("blogs", blog.Handle),
		}
		articles, err := NewArticleResource(restClient).List(blog.Id)
		if err != nil {
			return err
		}
		for _, article := range articles.Articles {
			contents.items = append(contents.items, Content{
				handle: article.Handle,
				html:   article.BodyHTML,
				meta:   articleToFrontMatter(article),
			})
		}
		blogs = append(blogs, contents)
	}

	all := append([]Contents{products, pages}, blogs...)
	errs := make([]error, len(all))
	var wg sync.WaitGroup
	progress := mpb.New(mpb.WithWaitGroup(&wg))
	for i, cts := range all {
		wg.Add(1)
		bar := progress.AddBar(int64(len(cts.items)),
			mpb.PrependDecorators(
				decor.Name(cts.kind),
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.AverageETA(decor.ET_STYLE_GO), "done",
				),
			),
		)

		go func(i int, o string, items []Content, b *mpb.Bar) {
			defer wg.Done()
			if errs[i] = dowloadContens(o, items, b); errs[i] != nil {
				// NOTE: Abort the bar so that progress.Wait doesn't wait for its completion.
				b.Abort(false)
			}
		}(i, path.Join(output, cts.kind), cts.items, bar)
	}

	progress.Wait()
	return errors.Join(errs...)
}
