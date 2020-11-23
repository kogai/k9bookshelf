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

	"github.com/Yamashou/gqlgenc/client"
	"github.com/gomarkdown/markdown"
	"github.com/mattn/godown"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

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

func download(output string) error {
	adminClient, ctx := gqlClient()
	res, err := adminClient.Products(ctx, 10)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Join(output, "products"), os.ModePerm)
	if err != nil {
		return err
	}

	// TODO: Use goroutine
	var downloadGroup sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&downloadGroup))
	bar := p.AddBar(int64(len(res.Products.Edges)),
		mpb.PrependDecorators(
			decor.Name(path.Join(output, "products")),
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
			),
		),
	)

	for _, edge := range res.Products.Edges {
		downloadGroup.Add(1)
		c := make(chan error)

		go func(handle, descriptionHTML string) {
			defer downloadGroup.Done()
			defer bar.Increment()

			file, err := os.Create(path.Join(output, "products", handle+".md"))
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
		}(edge.Node.Handle, edge.Node.DescriptionHTML)
		err = <-c
		if err != nil {
			return err
		}
	}

	p.Wait()
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
