package syncdata

import (
	"context"
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

func deployContents(ctx context.Context, input string, gqlClient *generated.Client, bar *mpb.Bar) error {
	files, err := ioutil.ReadDir(input)
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)
		c := make(chan error)
		filename := file.Name()

		go func(handle, pathToFile string) {
			defer wg.Done()
			defer bar.Increment()

			productByHandle, err := gqlClient.ProductByHandle(ctx, handle)
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

			res, err := gqlClient.Deploy(
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
			path.Join(input, filename),
		)

		err = <-c
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

// Deploy uploads contents to store
func Deploy(input string) error {
	files, err := ioutil.ReadDir(path.Join(input, "products"))
	if err != nil {
		return err
	}
	gqlClient, ctx := establishGqlClient()
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

			productByHandle, err := gqlClient.ProductByHandle(ctx, handle)
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

			res, err := gqlClient.Deploy(
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
