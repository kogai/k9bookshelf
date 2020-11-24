package syncdata

import (
	"fmt"
	"io/ioutil"
	"k9bookshelf/generated"
	"path"
	"path/filepath"
	"sync"

	"github.com/gomarkdown/markdown"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func deployProducts(contents []Content, bar *mpb.Bar) error {
	gqlClient, ctx := establishGqlClient()
	wg := sync.WaitGroup{}
	for _, content := range contents {
		wg.Add(1)
		c := make(chan error)

		go func(handle, html string) {
			defer wg.Done()
			fmt.Println("start", handle)
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
			c <- nil
		}(content.handle, content.html)

		err := <-c
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

// Deploy uploads contents to store
func Deploy(input string) error {
	rawProducts, err := ioutil.ReadDir(path.Join(input, "products"))
	if err != nil {
		return err
	}
	products := []Content{}
	for _, file := range rawProducts {
		filename := file.Name()
		handle := filename[0 : len(filename)-len(filepath.Ext(filename))]
		md, err := ioutil.ReadFile(path.Join(input, "products", filename))
		if err != nil {
			return err
		}
		html := string(markdown.ToHTML(md, nil, nil))
		products = append(products, Content{
			handle: handle,
			html:   html,
		})
	}

	// pages, err := ioutil.ReadDir(path.Join(input, "pages"))
	// if err != nil {
	// 	return err
	// }

	wg := sync.WaitGroup{}
	p := mpb.New(mpb.WithWaitGroup(&wg))
	bar := p.AddBar(int64(len(products)),
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

	err = deployProducts(products, bar)
	if err != nil {
		return err
	}

	p.Wait()
	return nil
}
